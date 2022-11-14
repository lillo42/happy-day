package reservation

import (
	"context"

	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	CreateRequest struct {
		Products []CreateProductRequest `json:"products"`
		Price    float64                `json:"price"`
		Discount float64                `json:"discount"`
	}

	CreateProductRequest struct {
		Id       uuid.UUID `json:"id"`
		Quantity int64     `json:"quantity"`
	}

	CreateHandler struct {
		productRepository     infrastructure.ProductRepository
		reservationRepository infrastructure.ReservationRepository
	}
)

func (handler CreateHandler) Handle(ctx context.Context, req CreateRequest) (reservation.State, error) {
	var state reservation.State
	if len(req.Products) == 0 {
		return state, infrastructure.ErrProductListIsEmpty
	}

	products, err := handler.productRepository.GetByProducts(ctx, common.Map(req.Products, func(item CreateProductRequest) uuid.UUID {
		return item.Id
	}))

	if err != nil {
		return state, err
	}

	state = reservation.State{
		Discount:   req.Discount,
		Price:      req.Price,
		FinalPrice: req.Price - req.Discount,
		Products: common.Map(products, func(item product.State) reservation.Product {
			p, _ := common.First(req.Products,
				func(item2 CreateProductRequest) bool {
					return item2.Id == item.Id
				},
				CreateProductRequest{})
			return reservation.Product{
				Id:       item.Id,
				Price:    item.Price,
				Quantity: p.Quantity,
			}
		}),
	}

	return handler.reservationRepository.Save(ctx, state)
}
