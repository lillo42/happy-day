package application

import (
	"context"
	"errors"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteProductRequest struct {
		Id uuid.UUID
	}

	DeleteProductHandler struct {
		repository infrastructure.ProductRepository
	}
)

var (
	ErrExistOtherProductWithThisProduct = errors.New("Exist other product with this product")
)

func (handler DeleteProductHandler) Handle(ctx context.Context, req DeleteProductRequest) error {
	exits, err := handler.repository.ExistAnyWithProduct(ctx, req.Id)
	if err != nil {
		return err
	}

	if exits {
		return ErrExistOtherProductWithThisProduct
	}

	return handler.repository.Delete(ctx, req.Id)
}
