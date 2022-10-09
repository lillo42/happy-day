package application

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuoteReservationWhenErrToGetComposed(t *testing.T) {
	req := QuoteReservationRequest{
		Products: []QuoteReservationProductRequest{
			{
				Id:     uuid.New(),
				Amount: 120,
			},
			{
				Id:     uuid.New(),
				Amount: 10,
			},
		},
	}

	expectedErr := errors.New(common.RandString(10))

	repository := &infrastructure.MockProductRepository{}
	repository.On("GetComposed", mock.Anything, mock.Anything).Return([]product.State{}, expectedErr)

	handler := QuoteReservationHandler{
		repository: repository,
	}

	_, err := handler.Handler(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestQuoteReservationWhenErrToGetProducts(t *testing.T) {
	req := QuoteReservationRequest{
		Products: []QuoteReservationProductRequest{
			{
				Id:     uuid.New(),
				Amount: 120,
			},
			{
				Id:     uuid.New(),
				Amount: 10,
			},
		},
	}

	expectedErr := errors.New(common.RandString(10))

	repository := &infrastructure.MockProductRepository{}
	repository.On("GetComposed", mock.Anything, mock.Anything).Return([]product.State{}, nil)

	repository.On("GetByProducts", mock.Anything, mock.Anything).Return([]product.State{}, expectedErr)

	handler := QuoteReservationHandler{
		repository: repository,
	}

	_, err := handler.Handler(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestQuoteReservation(t *testing.T) {
	type arg struct {
		products   []QuoteReservationProductRequest
		boxes      []product.State
		dbProducts []product.State
	}

	type param struct {
		name     string
		arg      arg
		expected float64
	}

	cases := []param{
		{
			name: "All product in one box",
			arg: arg{
				products: []QuoteReservationProductRequest{
					{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 40},
					{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 10},
				},
				boxes: []product.State{
					{
						Id:    uuid.New(),
						Price: 2.5,
						Products: []product.Product{
							{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 4},
							{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 1},
						},
					},
				},
				dbProducts: []product.State{},
			},
			expected: 25,
		},
		{
			name: "All product not in the box",
			arg: arg{
				products: []QuoteReservationProductRequest{
					{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 40},
					{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 10},
				},
				boxes: []product.State{
					{
						Id:    uuid.New(),
						Price: 2.5,
						Products: []product.Product{
							{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 4},
							{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 1},
							{Id: uuid.MustParse("b387b182-e12a-4757-99c8-31b6596d102d"), Amount: 1},
						},
					},
				},
				dbProducts: []product.State{
					{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Price: 1.5},
					{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Price: 4},
				},
			},
			expected: 100,
		},
		{
			name: "mixes",
			arg: arg{
				products: []QuoteReservationProductRequest{
					{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 45},
					{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 10},
					{Id: uuid.MustParse("b387b182-e12a-4757-99c8-31b6596d102d"), Amount: 3},
				},
				boxes: []product.State{
					{
						Id:    uuid.New(),
						Price: 2.5,
						Products: []product.Product{
							{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Amount: 4},
							{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Amount: 1},
						},
					},
				},
				dbProducts: []product.State{
					{Id: uuid.MustParse("1fd49fff-e5c6-4878-9c42-f2fa77bdc2fe"), Price: 1.5},
					{Id: uuid.MustParse("4491c392-4b6e-4a6e-aa64-c2897a258451"), Price: 4},
					{Id: uuid.MustParse("b387b182-e12a-4757-99c8-31b6596d102d"), Price: 2},
				},
			},
			expected: 38.5,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(k *testing.T) {

			req := QuoteReservationRequest{
				Products: c.arg.products,
			}

			repository := &infrastructure.MockProductRepository{}
			repository.On("GetComposed", mock.Anything, mock.Anything).Return(c.arg.boxes, nil)
			repository.On("GetByProducts", mock.Anything, mock.Anything).Return(c.arg.dbProducts, nil)

			handler := QuoteReservationHandler{
				repository: repository,
			}

			res, err := handler.Handler(context.Background(), req)
			assert.Nil(k, err)
			assert.Equal(k, c.expected, res.Price)
		})
	}
}
