package infrastructure

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"math"
	"time"

	"happy_day/domain/customer"

	"github.com/google/uuid"
)

const (
	CustomersCollection = "customers"

	CustomerIdAsc       CustomerSortBy = "id"
	CustomerIdDesc      CustomerSortBy = "id desc"
	CustomerNameAsc     CustomerSortBy = "name"
	CustomerNameDesc    CustomerSortBy = "name desc"
	CustomerCommentAsc  CustomerSortBy = "comment"
	CustomerCommentDesc CustomerSortBy = "comment desc"
)

var (
	_ CustomerRepository = (*GormCustomerRepository)(nil)

	ErrCustomerConcurrencyIssue = errors.New("product concurrency issue")
	ErrCustomerNotFound         = errors.New("product not found")
)

type (
	CustomerSortBy string
	CustomerFilter struct {
		Name    string
		Comment string
		Page    int64
		Size    int64
		SortBy  CustomerSortBy
	}

	CustomerRepository interface {
		GetById(ctx context.Context, id uint) (customer.State, error)
		GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error)

		Save(ctx context.Context, state customer.State) (customer.State, error)
		Delete(ctx context.Context, id uint) error
	}

	GormCustomerRepository struct {
		db *gorm.DB
	}
)

func (repository *GormCustomerRepository) GetById(ctx context.Context, id uint) (customer.State, error) {
	var customerDB Customer

	db := repository.db.
		WithContext(ctx).
		First(&customerDB, id)

	if db.Error != nil {
		err := db.Error
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			err = ErrCustomerNotFound
		}

		return customer.State{}, err
	}

	state := customer.State{
		ID:        customerDB.ID,
		Name:      customerDB.Name,
		Comment:   customerDB.Comment,
		Phones:    make([]customer.Phone, len(customerDB.Phones)),
		CreatedAt: customerDB.CreatedAt,
		UpdateAt:  customerDB.UpdatedAt,
	}

	for phoneIndex, phone := range customerDB.Phones {
		state.Phones[phoneIndex] = customer.Phone{
			Number: phone.Number,
		}
	}

	return state, nil
}

func (repository *GormCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error) {
	db := repository.db.WithContext(ctx)

	if len(filter.Name) > 0 {
		db.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if len(filter.Comment) > 0 {
		db.Where("comment LIKE ?", "%"+filter.Comment+"%")
	}

	var counter *int64
	db.Count(counter)

	var page Page[customer.State]
	if db.Error != nil {
		return page, db.Error
	}

	var customers []Customer
	db.
		Limit(int(filter.Size)).
		Offset(int((filter.Page - 1) * filter.Size)).
		Order(filter.SortBy).
		Preload("Phones").
		Find(&customers)

	if db.Error != nil {
		return page, db.Error
	}

	total := *counter
	page.TotalPages = page.TotalElements
	page.Items = make([]customer.State, len(customers))
	page.TotalElements = total
	if total > 0 {
		tmp := float64(total) / float64(filter.Size)
		tmp = math.Ceil(tmp)
		page.TotalPages = int64(tmp)
	}

	for index, customerDB := range customers {
		state := customer.State{
			ID:        customerDB.ID,
			Name:      customerDB.Name,
			Comment:   customerDB.Comment,
			Phones:    make([]customer.Phone, len(customerDB.Phones)),
			CreatedAt: customerDB.CreatedAt,
			UpdateAt:  customerDB.UpdatedAt,
		}

		for phoneIndex, phone := range customerDB.Phones {
			state.Phones[phoneIndex] = customer.Phone{
				Number: phone.Number,
			}
		}

		page.Items[index] = state
	}

	return page, nil
}

func (repository *GormCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	customerDB := &Customer{
		ID:        state.ID,
		Name:      state.Name,
		Comment:   state.Comment,
		CreatedAt: state.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err := repository.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			tx.WithContext(ctx)

			if customerDB.ID > 0 {
				tx.Where("updateAt = ?", state.UpdateAt)
			}

			tx.Save(customerDB)

			if tx.Error != nil {
				return tx.Error
			}

			if tx.RowsAffected == 0 {
				return ErrCustomerConcurrencyIssue
			}

			tx.
				WithContext(ctx).
				Where("customerID = ?", customerDB.ID).
				Delete(&CustomerPhone{})

			if tx.Error != nil {
				return tx.Error
			}

			for _, phone := range state.Phones {
				tx.Save(&CustomerPhone{
					ID:         uuid.New(),
					CustomerID: customerDB.ID,
					Number:     phone.Number,
				})

				if tx.Error != nil {
					return tx.Error
				}
			}

			return nil
		})

	if err != nil {
		return customer.State{}, err
	}

	return repository.GetById(ctx, customerDB.ID)
}

func (repository *GormCustomerRepository) Delete(ctx context.Context, id uint) error {
	db := repository.db.
		WithContext(ctx).
		Delete(&Customer{}, id)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return db.Error
}
