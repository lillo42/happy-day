package orders

import (
	"context"
	"github.com/google/uuid"
)

type (
	CustomerService interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}

	ProductService interface {
		Get(ctx context.Context, id uuid.UUID) (ProductProjection, error)
	}

	DiscountService interface {
		GetAll(ctx context.Context, productsID []uuid.UUID) ([]DiscountProjection, error)
	}

	DiscountProjection struct {
		Price    float64
		Products []DiscountProducts
	}

	DiscountProducts struct {
		ID       uuid.UUID
		Quantity uint
	}

	ProductProjection struct {
		ID    uuid.UUID
		Name  string
		Price float64
	}
)
