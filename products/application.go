package products

import (
	"context"
	"github.com/google/uuid"
)

var ()

type (
	Command struct {
		repository ProductRepository
	}

	CreateOrChange struct {
		ID       uuid.UUID    `json:"-"`
		Name     string       `json:"name"`
		Price    float64      `json:"price"`
		Products []BoxProduct `json:"products,omitempty"`
	}
)

func (c *Command) CreateOrChange(ctx context.Context, req CreateOrChange) (Product, error) {
	prod, err := c.repository.GetOrCreate(ctx, req.ID)

	if err != nil {
		return Product{}, err
	}

	if req.ID != uuid.Nil && prod.Version == 0 {
		return Product{}, ErrProductNotExists
	}

	if req.ID == uuid.Nil {
		prod.ID = uuid.New()
	}

	if len(req.Name) == 0 {
		return Product{}, ErrNameIsEmpty
	}

	if len(req.Name) > 255 {
		return Product{}, ErrNameTooLarge
	}

	if req.Price < 0 {
		return Product{}, ErrPriceIsInvalid
	}

	prod.Name = req.Name
	prod.Price = req.Price
	prod.Products = make([]BoxProduct, len(req.Products))
	for i, box := range req.Products {
		exists, err := c.repository.Exists(ctx, box.ID)
		if err != nil {
			return Product{}, err
		}

		if !exists {
			return Product{}, ErrBoxProductNotExists
		}

		prod.Products[i] = BoxProduct{
			ID:       box.ID,
			Quantity: box.Quantity,
		}
	}

	return c.repository.Save(ctx, prod)
}

func (c *Command) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repository.Delete(ctx, id)
}
