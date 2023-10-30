package discounts

import (
	"context"
	"github.com/google/uuid"
)

type (
	ProductService interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}
)
