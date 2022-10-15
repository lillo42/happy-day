package application

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrChangeCustomerHandlerWhenNameIsEmpty(t *testing.T) {
	handler := CreateOrChangeCustomerHandler{}
	_, err := handler.Handle(context.Background(), CreateOrChangeCustomerRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrCustomerNameIsEmpty, err)
}

func TestCreateOrChangeCustomerHandlerWhenPhonesIsEmpty(t *testing.T) {
	handler := CreateOrChangeCustomerHandler{}
	_, err := handler.Handle(context.Background(), CreateOrChangeCustomerRequest{
		State: customer.State{
			Name: common.RandString(10),
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrCustomerPhonesIsEmpty, err)
}

func TestCreateOrChangeCustomerHandlerWhenPhoneIsInvalid(t *testing.T) {
	phones := []customer.Phone{
		{
			Number: "1234567",
		},
		{
			Number: "1234567890123",
		},
		{
			Number: common.RandString(9),
		},
	}

	for _, phone := range phones {
		t.Run("Phone number"+phone.Number, func(t *testing.T) {
			handler := CreateOrChangeCustomerHandler{}
			_, err := handler.Handle(context.Background(), CreateOrChangeCustomerRequest{
				State: customer.State{
					Name: common.RandString(10),
					Phones: []customer.Phone{
						phone,
					},
				},
			})
			assert.NotNil(t, err)
			assert.Equal(t, ErrCustomerPhoneIsInvalid, err)
		})
	}
}

func TestCreateOrChangeCustomerHandlerWhenErrToGetById(t *testing.T) {
	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, expectedErr)
	handler := CreateOrChangeCustomerHandler{repository: repo}
	_, err := handler.Handle(context.Background(), CreateOrChangeCustomerRequest{
		State: customer.State{
			Id:   uuid.New(),
			Name: common.RandString(10),
			Phones: []customer.Phone{
				{
					Number: "123456789",
				},
			},
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreateOrChangeCustomerHandler(t *testing.T) {
	phones := []customer.Phone{
		{
			Number: "123456789",
		},
		{
			Number: "+1234567890",
		},
	}

	for _, phone := range phones {
		t.Run("Phone number"+phone.Number, func(t *testing.T) {
			repo := &infrastructure.MockCustomerRepository{}
			repo.
				On("Save", mock.Anything, mock.Anything).
				Return(customer.State{}, nil)
			handler := CreateOrChangeCustomerHandler{repository: repo}
			_, err := handler.Handle(context.Background(), CreateOrChangeCustomerRequest{
				State: customer.State{
					Name: common.RandString(10),
					Phones: []customer.Phone{
						phone,
					},
				},
			})
			assert.Nil(t, err)
		})
	}
}
