package main

import (
	"context"
	"github.com/google/uuid"
	"happyday/customers"
	"happyday/orders"
	"happyday/products"
)

type (
	GlobalProductService struct {
		repository products.ProductRepository
	}

	GlobalCustomerService struct {
		repository customers.CustomerRepository
	}
)

func (c *GlobalCustomerService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	cus, err := c.repository.GetOrCreate(ctx, id)
	if err != nil {
		return false, err
	}

	return cus.Version > 0, nil
}

func (s *GlobalProductService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	prod, err := s.repository.GetOrCreate(ctx, id)
	if err != nil {
		return false, err
	}

	return prod.Version > 0, nil
}

func (s *GlobalProductService) Get(ctx context.Context, id uuid.UUID) (orders.ProductProjection, error) {
	prod, err := s.repository.GetOrCreate(ctx, id)
	if err != nil {
		return orders.ProductProjection{}, err
	}

	return orders.ProductProjection{
		ID:    prod.ID,
		Name:  prod.Name,
		Price: prod.Price,
	}, err
}
