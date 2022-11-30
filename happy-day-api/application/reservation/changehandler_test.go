package reservation

import (
	"context"
	"errors"
	"testing"
	"time"

	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeReservationHandlerWhenPaymentInstallHasAmountLessThanZero(t *testing.T) {
	req := ChangeRequest{
		PaymentInstallments: []reservation.PaymentInstallment{
			{
				Method: reservation.Pix,
				Amount: -1,
				At:     time.Now().UTC(),
			},
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationPaymentInstallmentAmount, err)
}

func TestChangeReservationHandlerWhenCustomerIsInvalid(t *testing.T) {
	req := ChangeRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: "",
			},
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerNameIsEmpty, err)
}

func TestChangeReservationHandlerWhenAddressCityIsEmpty(t *testing.T) {
	req := ChangeRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationAddressCityIsEmpty, err)
}

func TestChangeReservationHandlerWhenAddressStreetIsEmpty(t *testing.T) {
	req := ChangeRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			City: common.RandString(10),
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationAddressStreetIsEmpty, err)
}

func TestChangeReservationHandlerWhenAddressNumberIsEmpty(t *testing.T) {
	req := ChangeRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			City:   common.RandString(10),
			Street: common.RandString(10),
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationAddressNumberIsInvalid, err)
}

func TestChangeReservationHandlerWhenAddressPostalCodeIsEmpty(t *testing.T) {
	req := ChangeRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			City:   common.RandString(10),
			Street: common.RandString(10),
			Number: common.RandString(10),
		},
	}

	handler := ChangeHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationAddressPostalCodeIsEmpty, err)
}

func TestChangeReservationHandlerWhenGetById(t *testing.T) {
	req := ChangeRequest{
		Id: uuid.New(),
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			City:       common.RandString(10),
			Street:     common.RandString(10),
			Number:     common.RandString(10),
			PostalCode: common.RandString(10),
		},
	}

	expectedErr := errors.New(common.RandString(10))
	repository := &infrastructure.MockReservationRepository{}
	repository.
		On("GetById", mock.Anything, req.Id).
		Return(reservation.State{}, expectedErr)

	handler := ChangeHandler{repository: repository}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestChangeReservationHandler(t *testing.T) {
	req := ChangeRequest{
		Id: uuid.New(),
		Customer: reservation.Customer{
			State: customer.State{
				Name: common.RandString(10),
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			City:       common.RandString(10),
			Street:     common.RandString(10),
			Number:     common.RandString(10),
			PostalCode: common.RandString(10),
		},
	}

	repository := &infrastructure.MockReservationRepository{}
	repository.
		On("GetById", mock.Anything, req.Id).
		Return(reservation.State{}, nil)

	repository.
		On("Save", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	handler := ChangeHandler{repository: repository}
	_, err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
