package application

import (
	"context"
	"happy_day/domain/product"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type GetProductByIdHandler struct {
	repository infrastructure.ProductRepository
}

func (handler GetProductByIdHandler) Handle(ctx context.Context, req uuid.UUID) (product.State, error) {
	return handler.repository.GetById(ctx, req)
}
