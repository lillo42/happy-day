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
		Product  []CreateReservationProductRequest `json:"product"`
		Price    float64                           `json:"price"`
		Discount float64                           `json:"discount"`
	}

	CreateReservationProductRequest struct {
		Id     uuid.UUID `json:"id"`
		Amount int64     `json:"amount"`
	}

	CreateReservationHandler struct {
		productRepository     infrastructure.ProductRepository
		reservationRepository infrastructure.ReservationRepository
	}
)

var (
	ErrProductListIsEmpty = errors.New("Product list cannot be empty")
)

func (handler CreateReservationHandler) Handle(ctx context.Context, req CreateReservationRequest) (reservation.State, error) {

	var state reservation.State

	if len(req.Product) == 0 {
		return state, ErrProductListIsEmpty
	}

	products, err := handler.productRepository.GetByProducts(ctx, common.Map(req.Product, func(item CreateReservationProductRequest) uuid.UUID {
		return item.Id
	}))

	if err != nil {
		return state, err
	}

	state = reservation.State{
		Id:       uuid.New(),
		Discount: req.Discount,
		Price:    req.Price,
		Products: common.MapWithIndex(products, func(index int, item product.State) reservation.Product {
			return reservation.Product{
				Id:     item.Id,
				Price:  item.Price,
				Amount: req.Product[index].Amount,
			}
		}),
	}

	return handler.reservationRepository.Save(ctx, state)
}
