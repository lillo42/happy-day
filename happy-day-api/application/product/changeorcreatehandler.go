package product

import (
	"context"

	"happy_day/domain/product"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	ChangeOrCreateRequest struct {
		product.State
	}

	ChangeOrCreateHandler struct {
		repository infrastructure.ProductRepository
	}
)

func (handler ChangeOrCreateHandler) Handle(ctx context.Context, req ChangeOrCreateRequest) (product.State, error) {
	if len(req.Name) == 0 {
		return req.State, infrastructure.ErrProductNameIsEmpty
	}

	if req.Price < 0 {
		return req.State, infrastructure.ErrProductPriceIsLessThanZero
	}

	for _, p := range req.Products {
		exists, err := handler.repository.Exists(ctx, p.Id)
		if err != nil {
			return req.State, err
		}

		if !exists {
			return req.State, infrastructure.ErrProductNotFound
		}

		if p.Quantity <= 0 {
			return req.State, infrastructure.ErrProductAmountIsInvalid
		}
	}

	if req.Id != uuid.Nil {
		state, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return state, err
		}

		req.CreatedAt = state.CreatedAt
		req.ModifiedAt = state.ModifiedAt
	}

	return handler.repository.Save(ctx, req.State)
}
