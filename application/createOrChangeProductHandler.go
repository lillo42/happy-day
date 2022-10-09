package application

import (
	"context"
	"happy_day/domain/product"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	CreateOrChangeProductRequest struct {
		product.State
	}

	CreateOrChangeProductHandler struct {
		repository infrastructure.ProductRepository
	}
)

func (handler CreateOrChangeProductHandler) Handler(ctx context.Context, req CreateOrChangeProductRequest) (product.State, error) {
	product := product.State{}
	if req.Id != uuid.Nil {
		product, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return product, err
		}
	}

	return product, nil
}
