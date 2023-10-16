package customer

import (
	"context"
	"errors"
	"testing"

	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrChangeCustomerHandlerWhenNameIsEmpty(t *testing.T) {
	handler := ChangeOrCreateHandler{}
	_, err := handler.Handle(context.Background(), ChangeOrCreateRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerNameIsEmpty, err)
}

func TestCreateOrChangeCustomerHandlerWhenPhonesIsEmpty(t *testing.T) {
	handler := ChangeOrCreateHandler{}
	_, err := handler.Handle(context.Background(), ChangeOrCreateRequest{
		State: customer.State{
			Name: common.RandString(10),
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerPhonesIsEmpty, err)
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
			handler := ChangeOrCreateHandler{}
			_, err := handler.Handle(context.Background(), ChangeOrCreateRequest{
				State: customer.State{
					Name: common.RandString(10),
					Phones: []customer.Phone{
						phone,
					},
				},
			})
			assert.NotNil(t, err)
			assert.Equal(t, infrastructure.ErrCustomerPhoneIsInvalid, err)
		})
	}
}

func TestCreateOrChangeCustomerHandlerWhenErrToGetById(t *testing.T) {
	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, expectedErr)
	handler := ChangeOrCreateHandler{repository: repo}
	_, err := handler.Handle(context.Background(), ChangeOrCreateRequest{
		State: customer.State{
			ID:   1,
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
			handler := ChangeOrCreateHandler{repository: repo}
			_, err := handler.Handle(context.Background(), ChangeOrCreateRequest{
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
