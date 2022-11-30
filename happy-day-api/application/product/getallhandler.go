package product

import (
	"context"

	"happy_day/domain/product"
	"happy_day/infrastructure"
)

type GetAllHandler struct {
	repository infrastructure.ProductRepository
}

func (handler GetAllHandler) Handle(ctx context.Context, req infrastructure.ProductFilter) (infrastructure.Page[product.State], error) {
	if req.Size <= 0 {
		req.Size = 50
	}

	if req.Page <= 1 {
		req.Page = 1
	}

	return handler.repository.GetAll(ctx, req)
}
