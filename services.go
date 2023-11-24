package main

import (
	"context"
	"github.com/google/uuid"
	"happyday/customers"
	"happyday/discounts"
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

	GlobalDiscountService struct {
		repository discounts.DiscountRepository
	}
)

func (g *GlobalDiscountService) GetAll(ctx context.Context, productsID []uuid.UUID) ([]orders.DiscountProjection, error) {
	discounts, err := g.repository.GetAllWithProducts(ctx, productsID)
	if err != nil {
		return nil, err
	}

	proj := make([]orders.DiscountProjection, len(discounts))
	for i, discount := range discounts {
		prods := make([]orders.DiscountProducts, len(discount.Products))
		for j, prod := range discount.Products {
			prods[j] = orders.DiscountProducts{
				ID:       prod.ID,
				Quantity: prod.Quantity,
			}
		}

		proj[i] = orders.DiscountProjection{
			Price:    discount.Price,
			Products: prods,
		}
	}

	return proj, nil
}

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
