package application

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReservationHandlerWhenProductListIsEmpty(t *testing.T) {
	handler := CreateReservationHandler{}
	_, err := handler.Handle(context.Background(), CreateReservationRequest{})

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductListIsEmpty, err)
}

func TestCreateReservationHandlerWhenErrToGetByProducts(t *testing.T) {
	expectedErr := errors.New(common.RandString(10))
	productRepository := &infrastructure.MockProductRepository{}
	productRepository.
		On("GetByProducts", mock.Anything, mock.Anything).
		Return([]product.State{}, expectedErr)

	handler := CreateReservationHandler{productRepository: productRepository}
	_, err := handler.Handle(context.Background(), CreateReservationRequest{
		Products: []CreateReservationProductRequest{
			{
				Id:       uuid.New(),
				Quantity: 2,
			},
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreateReservationHandler(t *testing.T) {
	req := CreateReservationRequest{
		Products: []CreateReservationProductRequest{
			{
				Id:       uuid.New(),
				Quantity: 4,
			},
			{
				Id:       uuid.New(),
				Quantity: 1,
			},
		},
	}

	productRepository := &infrastructure.MockProductRepository{}
	productRepository.
		On("GetByProducts", mock.Anything, mock.Anything).
		Return(common.Map(req.Products, func(item CreateReservationProductRequest) product.State {
			return product.State{
				Id:    item.Id,
				Price: float64(item.Quantity) * 1.5,
			}
		}), nil)

	reservationRepository := &infrastructure.MockReservationRepository{}
	reservationRepository.
		On("Save", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	handler := CreateReservationHandler{
		productRepository:     productRepository,
		reservationRepository: reservationRepository,
	}
	_, err := handler.Handle(context.Background(), req)

	assert.Nil(t, err)
}
