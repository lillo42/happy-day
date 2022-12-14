package reservation

import (
	"context"
	"math"

	"happy_day/common"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	QuoteRequest struct {
		Products []QuoteProductRequest `json:"products"`
	}

	QuoteProductRequest struct {
		Id       uuid.UUID `json:"id"`
		Quantity int64     `json:"quantity"`
	}

	QuoteProductResponse struct {
		Price float64 `json:"price"`
	}

	QuoteHandler struct {
		repository infrastructure.ProductRepository
	}
)

func (handler QuoteHandler) Handler(ctx context.Context, req QuoteRequest) (QuoteProductResponse, error) {
	ids := common.Map(req.Products, func(item QuoteProductRequest) uuid.UUID {
		return item.Id
	})

	composed, err := handler.repository.GetComposed(ctx, ids)
	if err != nil {
		return QuoteProductResponse{}, err
	}

	productAmount := map[uuid.UUID]int64{}
	for _, product := range req.Products {
		productAmount[product.Id] = product.Quantity
	}

	var price float64
	for _, box := range composed {
		min := math.MaxFloat64
		for _, product := range box.Products {
			if product.Quantity == 0 {
				break
			}

			amount, exists := productAmount[product.Id]
			if !exists {
				min = 0
				break
			}
			min = math.Min(min, math.Floor(float64(amount/product.Quantity)))
		}

		if min == 0 || min == math.MaxFloat64 {
			continue
		}

		price += min * box.Price
		for _, product := range box.Products {
			productAmount[product.Id] -= product.Quantity * int64(min)
		}
	}

	products, err := handler.repository.GetByProducts(ctx, ids)
	if err != nil {
		return QuoteProductResponse{}, err
	}

	for _, product := range products {
		price += product.Price * float64(productAmount[product.Id])
	}

	return QuoteProductResponse{
		Price: price,
	}, nil
}
