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
		Exists(ctx context.Context, id uuid.UUID) (bool, error)

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
		Items:      make([]Product, 0),
		TotalItems: counter,
		TotalPage:  totalPage,
	}

	for _, prodDB := range productsDB {
		var boxes []infra.BoxProduct
		result := g.db.
			WithContext(ctx).
			Table("box_products").
			Where("parent_id = ?", prodDB.ID).
			Scan(&boxes)

		if result.Error != nil {
			return infra.Page[Product]{}, result.Error
		}

		if boxes != nil {
			for i, box := range boxes {
				var tmp infra.Product
				result := g.db.
					WithContext(ctx).
					First(&tmp, box.ProductID)

				if result.Error != nil {
					return infra.Page[Product]{}, result.Error
				}

				boxes[i].Product = tmp
			}
		}

		if boxes == nil {
			boxes = make([]infra.BoxProduct, 0)
		}

		prodDB.Products = boxes
		page.Items = append(page.Items, mapToProduct(prodDB))
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

	var boxes []infra.BoxProduct
	result = g.db.
		WithContext(ctx).
		Table("box_products").
		Where("parent_id = ?", productDB.ID).
		Scan(&boxes)

	if result.Error != nil {
		return Product{}, result.Error
	}

	if boxes != nil {
		for i, box := range boxes {
			var tmp infra.Product
			result := g.db.
				WithContext(ctx).
				First(&tmp, box.ProductID)

			if result.Error != nil {
				return Product{}, result.Error
			}

			boxes[i].Product = tmp
		}
	}

	if boxes == nil {
		boxes = make([]infra.BoxProduct, 0)
	}

	productDB.Products = boxes
	return mapToProduct(productDB), nil
}

func (g *GormProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var productDB infra.Product
	result := g.db.
		WithContext(ctx).
		First(&productDB, "external_id = ?", id)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}

	return !errors.Is(result.Error, gorm.ErrRecordNotFound), nil
}

func (g *GormProductRepository) Save(ctx context.Context, product Product) (Product, error) {
	var productDB infra.Product
	if product.Version > 0 {
		result := g.db.
			WithContext(ctx).
			First(&productDB, "external_id = ?", product.ID)

		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return Product{}, result.Error
			} else {
				return Product{}, ErrProductNotExists
			}
		}

	} else {
		productDB.ExternalID = product.ID
		productDB.CreateAt = time.Now()
	}

	productDB.Name = product.Name
	productDB.Price = product.Price
	productDB.Products = make([]infra.BoxProduct, len(product.Products))
	productDB.Version = product.Version + 1
	productDB.UpdateAt = time.Now()

	for i, box := range product.Products {
		var boxDB infra.Product
		result := g.db.
			WithContext(ctx).
			First(&boxDB, "external_id = ?", box.ID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return Product{}, ErrProductNotExists
			} else {
				return Product{}, result.Error
			}
		}

		productDB.Products[i] = infra.BoxProduct{
			Quantity:  box.Quantity,
			Product:   boxDB,
			ProductID: boxDB.ID,
		}
	}

	err := g.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			result := tx.Where("version = ?", product.Version).Save(&productDB)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return ErrConcurrencyUpdate
			}

			result = tx.Delete(&infra.BoxProduct{}, "parent_id = ?", productDB.ID)
			if result.Error != nil {
				return result.Error
			}

			for _, box := range productDB.Products {
				box.Parent = productDB
				box.ParentID = productDB.ID

				result = tx.Save(&box)
				if result.Error != nil {
					return result.Error
				}
			}

			return nil
		})

	if err != nil {
		return Product{}, err
	}

	return mapToProduct(productDB), nil
}

func (g *GormProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var productDB infra.Product
	result := g.db.
		WithContext(ctx).
		First(&productDB, "external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		return result.Error
	}

	return g.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			result := tx.Delete(&infra.BoxProduct{}, "parent_id = ?", productDB.ID)
			if result.Error != nil {
				return result.Error
			}

			result = tx.Delete(&infra.BoxProduct{}, "product_id = ?", productDB.ID)
			if result.Error != nil {
				return result.Error
			}

			result = tx.Delete(&infra.Product{}, "id = ?", productDB.ID)

			if result.Error != nil {
				return result.Error
			}

			return nil
		})
}

func mapToProduct(productDB infra.Product) Product {
	product := Product{
		ID:       productDB.ExternalID,
		Name:     productDB.Name,
		Price:    productDB.Price,
		CreateAt: productDB.CreateAt,
		UpdateAt: productDB.UpdateAt,
		Version:  productDB.Version,
		Products: make([]BoxProduct, len(productDB.Products)),
	}

	for i, box := range productDB.Products {
		product.Products[i] = BoxProduct{
			ID:       box.Product.ExternalID,
			Name:     box.Product.Name,
			Quantity: box.Quantity,
		}
	}

	return product

}
