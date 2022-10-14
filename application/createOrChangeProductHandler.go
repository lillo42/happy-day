package application

import (
	"context"
	"errors"
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

var (
	ErrProductNameIsEmpty         = errors.New("product name is empty")
	ErrProductPriceIsLessThanZero = errors.New("product price is less than zero")
	ErrProductAmountIsInvalid     = errors.New("product amount is invalid")
)

func (handler CreateOrChangeProductHandler) Handle(ctx context.Context, req CreateOrChangeProductRequest) (product.State, error) {
	if len(req.Name) == 0 {
		return req.State, ErrProductNameIsEmpty
	}

	if req.Price < 0 {
		return req.State, ErrProductPriceIsLessThanZero
	}

	for _, product := range req.Products {
		exists, err := handler.repository.Exists(ctx, product.Id)
		if err != nil {
			return req.State, err
		}

		if !exists {
			return req.State, infrastructure.ErrProductNotFound
		}

		if product.Amount <= 0 {
			return req.State, ErrProductAmountIsInvalid
		}
	}

	if req.Id != uuid.Nil {
		state, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return state, err
		}

		req.CreateAt = state.CreateAt
		req.ModifyAt = state.ModifyAt
	}

	return handler.repository.Save(ctx, req.State)
}
