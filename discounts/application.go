package discounts

import (
	"context"
	"github.com/google/uuid"
)

type (
	Command struct {
		repository     DiscountRepository
		productService ProductService
	}

	CreateOrChange struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Price    float64   `json:"price"`
		Products []Product `json:"products"`
	}
)

func (c *Command) CreateOrChange(ctx context.Context, req CreateOrChange) (Discount, error) {
	discount, err := c.repository.GetOrCreate(ctx, req.ID)

	if err != nil {
		return Discount{}, err
	}

	if req.ID != uuid.Nil && discount.Version == 0 {
		return Discount{}, ErrDiscountNotFound
	}

	if req.ID == uuid.Nil {
		discount.ID = uuid.New()
	}

	if len(req.Name) == 0 {
		return Discount{}, ErrNameIsEmpty
	}

	if len(req.Name) > 255 {
		return Discount{}, ErrNameIsTooLarge
	}

	if req.Price <= 0 {
		return Discount{}, ErrPriceIsInvalid
	}

	if req.Products == nil || len(req.Products) == 0 {
		return Discount{}, ErrProductsIsMissing
	}

	discount.Name = req.Name
	discount.Price = req.Price
	discount.Products = make([]Product, len(req.Products))

	for i, prod := range req.Products {
		exists, err := c.productService.Exists(ctx, prod.ID)
		if err != nil {
			return Discount{}, err
		}

		if !exists {
			return Discount{}, ErrProductNotFound
		}

		discount.Products[i] = Product{
			ID:       prod.ID,
			Quantity: prod.Quantity,
		}
	}

	return c.repository.Save(ctx, discount)
}

func (c *Command) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repository.Delete(ctx, id)
}
