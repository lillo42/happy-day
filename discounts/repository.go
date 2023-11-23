package discounts

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"happyday/infra"
	"time"
)

var _ DiscountRepository = (*GormDiscountRepository)(nil)

type (
	DiscountRepository interface {
		GetOrCreate(ctx context.Context, id uuid.UUID) (Discount, error)
		GetAll(ctx context.Context, filter DiscountFilter) (infra.Page[Discount], error)

		Save(ctx context.Context, discount Discount) (Discount, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	GormDiscountRepository struct {
		db *gorm.DB
	}

	DiscountFilter struct {
		Name string
		Page int
		Size int
	}
)

func (g *GormDiscountRepository) GetAll(ctx context.Context, filter DiscountFilter) (infra.Page[Discount], error) {
	query := g.db.
		WithContext(ctx).
		Model(&infra.Discount{})

	if len(filter.Name) > 0 {
		query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	var counter int64
	result := query.Count(&counter)
	if result.Error != nil {
		return infra.Page[Discount]{}, result.Error
	}

	var discountsDB []infra.Discount
	result = query.
		Preload("Products").
		Preload("Products.Product").
		Limit(filter.Size).
		Offset(filter.Page * filter.Size).
		Scan(&discountsDB)

	if result.Error != nil {
		return infra.Page[Discount]{}, result.Error
	}

	var totalPage int64
	if counter > 0 {
		totalPage = counter / int64(filter.Size)
	}

	page := infra.Page[Discount]{
		Items:      make([]Discount, len(discountsDB)),
		TotalItems: counter,
		TotalPage:  totalPage,
	}

	for i, discountDB := range discountsDB {
		page.Items[i] = mapToDiscount(discountDB)
	}

	return page, nil
}

func (g *GormDiscountRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Discount, error) {
	if id == uuid.Nil {
		return Discount{}, nil
	}

	var discount infra.Discount
	result := g.db.
		WithContext(ctx).
		Preload("Products").
		Preload("Products.Product").
		First(&discount, "external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Discount{ID: id}, nil
		}
		return Discount{}, result.Error
	}

	return mapToDiscount(discount), nil
}

func (g *GormDiscountRepository) Save(ctx context.Context, discount Discount) (Discount, error) {
	var discountDB infra.Discount
	if discount.Version == 0 {
		discountDB.ExternalID = discount.ID
		discountDB.CreateAt = time.Now()
	} else {
		result := g.db.
			WithContext(ctx).
			First(&discountDB, "external_id = ?", discount.ID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return Discount{}, ErrDiscountNotFound
			}
			return Discount{}, result.Error
		}
	}

	discountDB.Products = make([]infra.DiscountProducts, len(discount.Products))
	for i, prod := range discount.Products {
		var prodDB infra.Product

		result := g.db.
			WithContext(ctx).
			First(&prodDB, "external_id = ?", prod.ID)

		if result.Error != nil {
			return Discount{}, result.Error
		}

		discountDB.Products[i] = infra.DiscountProducts{
			ProductID: prodDB.ID,
			Quantity:  prod.Quantity,
		}
	}

	discountDB.Name = discount.Name
	discountDB.Price = discount.Price
	discountDB.Version = discount.Version + 1
	discountDB.UpdateAt = time.Now()

	err := g.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if discountDB.ID > 0 {
				result := tx.
					Where("discount_id = ?", discountDB.ID).
					Delete(&infra.DiscountProducts{})

				if result.Error != nil {
					return result.Error
				}
			}

			result := tx.
				Where("version = ?", discount.Version).
				Save(&discountDB)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return ErrConcurrencyUpdate
			}

			return nil
		})

	if err != nil {
		return Discount{}, err
	}

	return mapToDiscount(discountDB), nil
}

func (g *GormDiscountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var discount infra.Discount
	result := g.db.
		WithContext(ctx).
		First(&discount, "external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		return result.Error
	}

	return g.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			result := tx.
				WithContext(ctx).
				Where("discount_id = ?", discount.ID).
				Delete(&infra.DiscountProducts{})

			if result.Error != nil {
				return result.Error
			}

			result = tx.
				WithContext(ctx).
				Delete(&discount)

			return result.Error
		})
}

func mapToDiscount(discountDB infra.Discount) Discount {
	discount := Discount{
		ID:       discountDB.ExternalID,
		Name:     discountDB.Name,
		Price:    discountDB.Price,
		Version:  discountDB.Version,
		CreateAt: discountDB.CreateAt,
		UpdateAt: discountDB.UpdateAt,
	}

	var length int
	if discountDB.Products != nil {
		length = len(discountDB.Products)
	}

	discount.Products = make([]Product, length)

	for i, prod := range discountDB.Products {
		discount.Products[i] = Product{
			ID:       prod.Product.ExternalID,
			Name:     prod.Product.Name,
			Quantity: prod.Quantity,
		}
	}

	return discount
}
