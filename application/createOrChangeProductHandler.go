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
	ErrProductNameIsEmpty         = errors.New("Product name is empty")
	ErrProductPriceIsLessThanZero = errors.New("Product price is less than zero")
	ErrProductNotFound            = errors.New("Product not found")
	ErrProductAmoutIsInvalid      = errors.New("Product amount is invalid")
)

func (handler CreateOrChangeProductHandler) Handle(ctx context.Context, req CreateOrChangeProductRequest) (product.State, error) {
	state := req.State
	if req.Id != uuid.Nil {
		product, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return product, err
		}
	}

	if len(req.Name) == 0 {
		return state, ErrProductNameIsEmpty
	}

	if req.Price < 0 {
		return state, ErrProductPriceIsLessThanZero
	}

	for _, product := range req.Products {
		exists, err := handler.repository.Exists(ctx, product.Id)
		if err != nil {
			return state, err
		}

		if !exists {
			return state, ErrProductNotFound
		}

		if product.Amount <= 0 {
			return state, ErrProductAmoutIsInvalid
		}
	}

	return handler.repository.Save(ctx, state)
}
