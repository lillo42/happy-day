package product

import (
	"context"

	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteRequest struct {
		Id uuid.UUID
	}

	DeleteHandler struct {
		repository infrastructure.ProductRepository
	}
)

func (handler DeleteHandler) Handle(ctx context.Context, req DeleteRequest) error {
	exits, err := handler.repository.ExistAnyWithProduct(ctx, req.Id)
	if err != nil {
		return err
	}

	if exits {
		return infrastructure.ErrExistOtherProductWithThisProduct
	}

	return handler.repository.Delete(ctx, req.Id)
}
