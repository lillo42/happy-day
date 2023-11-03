package products

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"happyday/infra"
	"time"
)

type (
	ProductFilter struct {
		Name string
		Page int
		Size int
	}

	ProductRepository interface {
		GetAll(ctx context.Context, filter ProductFilter) (infra.Page[Product], error)
		GetOrCreate(ctx context.Context, id uuid.UUID) (Product, error)

		Save(ctx context.Context, product Product) (Product, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	GormProductRepository struct {
		db *gorm.DB
	}
)

func (g *GormProductRepository) GetAll(ctx context.Context, filter ProductFilter) (infra.Page[Product], error) {
	query := g.db.
		WithContext(ctx).
		Table("products")

	if len(filter.Name) > 0 {
		query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	var counter int64
	result := query.Count(&counter)
	if result.Error != nil {
		return infra.Page[Product]{}, result.Error
	}

	var productsDB []infra.Product
	result = query.
		Order("id").
		Limit(filter.Size).
		Offset(filter.Page * filter.Size).
		Scan(&productsDB)

	if result.Error != nil {
		return infra.Page[Product]{}, result.Error
	}

	var totalPage int64
	if counter > 0 {
		totalPage = counter / int64(filter.Size)
	}

	page := infra.Page[Product]{
		Items:      make([]Product, len(productsDB)),
		TotalItems: counter,
		TotalPage:  totalPage,
	}

	for i, prodDB := range productsDB {
		page.Items[i] = mapToProduct(prodDB)
	}

	return page, nil
}

func (g *GormProductRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Product, error) {
	var productDB infra.Product
	result := g.db.
		WithContext(ctx).
		First(&productDB, "external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Product{ID: id}, nil
		}

		return Product{}, result.Error
	}

	return mapToProduct(productDB), nil
}

func (g *GormProductRepository) Save(ctx context.Context, product Product) (Product, error) {
	var productDB infra.Product
	if product.Version == 0 {
		productDB.ExternalID = uuid.New()
		productDB.CreateAt = time.Now()
	} else {
		result := g.db.
			WithContext(ctx).
			First(&productDB, "external_id = ?", product.ID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return Product{}, ErrProductNotExists
			}
			return Product{}, result.Error
		}
	}

	productDB.Name = product.Name
	productDB.Price = product.Price
	productDB.Version = product.Version + 1
	productDB.UpdateAt = time.Now()

	var result *gorm.DB
	if product.Version > 1 {
		result = g.db.
			WithContext(ctx).
			Where("version = ?", product.Version).
			Save(&productDB)
	} else {
		result = g.db.
			WithContext(ctx).
			Save(&productDB)
	}

	if result.Error != nil {
		return Product{}, result.Error
	}

	return mapToProduct(productDB), nil
}

func (g *GormProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var productDB infra.Product
	result := g.db.
		WithContext(ctx).
		Delete(&productDB, "external_id = ?", id)

	return result.Error
}

func mapToProduct(productDB infra.Product) Product {
	return Product{
		ID:       productDB.ExternalID,
		Name:     productDB.Name,
		Price:    productDB.Price,
		CreateAt: productDB.CreateAt,
		UpdateAt: productDB.UpdateAt,
		Version:  productDB.Version,
	}
}
