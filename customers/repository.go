package customers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"happyday/infra"
	"time"
)

var _ CustomerRepository = (*GormCustomerRepository)(nil)

type (
	CustomerFilter struct {
		Name    string
		Comment string
		Phone   string
		Page    int
		Size    int
	}

	CustomerRepository interface {
		GetAll(ctx context.Context, filter CustomerFilter) (infra.Page[Customer], error)
		GetOrCreate(ctx context.Context, id uuid.UUID) (Customer, error)
		Save(ctx context.Context, customer Customer) (Customer, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	GormCustomerRepository struct {
		db *gorm.DB
	}
)

func (g *GormCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (infra.Page[Customer], error) {
	query := g.db.WithContext(ctx).Table("customers")

	if len(filter.Name) > 0 {
		query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if len(filter.Comment) > 0 {
		query.Where("comment LIKE ?", "%"+filter.Comment+"%")
	}

	if len(filter.Phone) > 0 {
		query.Where("phones LIKE ?", "%"+filter.Phone+"%")
	}

	var counter int64
	result := query.Count(&counter)
	if result.Error != nil || counter == 0 {
		return infra.Page[Customer]{}, result.Error
	}

	var customerDBs []infra.Customer
	result = query.
		Limit(filter.Size).
		Offset(filter.Page * filter.Size).
		Scan(&customerDBs)
	if result.Error != nil {
		return infra.Page[Customer]{}, result.Error
	}

	var totalPage int64
	if counter > 0 {
		totalPage = counter / int64(filter.Size)
	}

	var page infra.Page[Customer]
	page.TotalItems = counter
	page.TotalPage = totalPage

	var items []Customer
	for _, customerDB := range customerDBs {
		customer, err := mapToCustomer(customerDB)

		if err != nil {
			return infra.Page[Customer]{}, err
		}

		items = append(items, customer)
	}

	return page, nil
}

func (g *GormCustomerRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Customer, error) {
	var customerDB infra.Customer
	db := g.db.
		WithContext(ctx).
		First(&customerDB, "external_id = ?", id)

	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return Customer{}, db.Error
	}

	customerDB.ExternalID = id
	return mapToCustomer(customerDB)
}

func (g *GormCustomerRepository) Save(ctx context.Context, customer Customer) (Customer, error) {
	phones, _ := json.Marshal(customer.Phones)
	var customerDB infra.Customer
	if customer.Version > 0 {
		result := g.db.WithContext(ctx).First(&customerDB, "external_id = ?", customer.ID)
		if result.Error != nil {
			return Customer{}, result.Error
		}
	} else {
		customerDB.ExternalID = customer.ID
		customerDB.CreateAt = time.Now()
	}

	customerDB.Name = customer.Name
	customerDB.Comment = customer.Comment
	customerDB.Phones = string(phones)
	customerDB.Pix = customer.Pix
	customerDB.Version = customer.Version + 1
	customerDB.CreateAt = customer.CreateAt
	customerDB.UpdateAt = time.Now()

	query := g.db.
		WithContext(ctx)

	if customer.Version > 0 {
		query.Where("version = ?", customer.Version)
	}

	result := query.Save(&customerDB)

	if result.Error != nil {
		return Customer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return Customer{}, ErrConcurrencyUpdate
	}

	customer.UpdateAt = customerDB.UpdateAt
	customer.Version = customerDB.Version
	return customer, nil

}

func (g *GormCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := g.db.
		WithContext(ctx).
		Delete(&infra.Customer{}, "external_id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return result.Error
}

func mapToCustomer(customerDB infra.Customer) (Customer, error) {
	var phones []string
	if len(customerDB.Phones) > 0 {
		_ = json.Unmarshal([]byte(customerDB.Phones), &phones)
	}

	return Customer{
		ID:       customerDB.ExternalID,
		Name:     customerDB.Name,
		Comment:  customerDB.Comment,
		Phones:   phones,
		CreateAt: customerDB.CreateAt,
		UpdateAt: customerDB.UpdateAt,
		Version:  customerDB.Version,
	}, nil
}
