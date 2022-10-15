package application

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	CreateReservationRequest struct {
		Products []CreateReservationProductRequest `json:"products"`
		Price    float64                           `json:"price"`
		Discount float64                           `json:"discount"`
	}

	CreateReservationProductRequest struct {
		Id       uuid.UUID `json:"id"`
		Quantity int64     `json:"quantity"`
	}

	CreateReservationHandler struct {
		productRepository     infrastructure.ProductRepository
		reservationRepository infrastructure.ReservationRepository
	}
)

var (
	ErrProductListIsEmpty = errors.New("product list cannot be empty")
)

func (handler CreateReservationHandler) Handle(ctx context.Context, req CreateReservationRequest) (reservation.State, error) {
	var state reservation.State
	if len(req.Products) == 0 {
		return state, ErrProductListIsEmpty
	}

	products, err := handler.productRepository.GetByProducts(ctx, common.Map(req.Products, func(item CreateReservationProductRequest) uuid.UUID {
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
				func(item2 CreateReservationProductRequest) bool {
					return item2.Id == item.Id
				},
				CreateReservationProductRequest{})
			return reservation.Product{
				Id:     item.Id,
				Price:  item.Price,
				Amount: p.Quantity,
			}
		}),
	}

	return handler.reservationRepository.Save(ctx, state)
}
