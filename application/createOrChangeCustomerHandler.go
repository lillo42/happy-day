package application

import (
	"errors"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	ErrCustomerNameIsEmpty    = errors.New("customer name is empty")
	ErrCustomerPhonesIsEmpty  = errors.New("customer phones is empty")
	ErrCustomerPhoneIsInvalid = errors.New("customer phone is invalid")
)

type (
	CreateOrChangeCustomerRequest struct {
		customer.State
	}

	CreateOrChangeCustomerHandler struct {
		repository infrastructure.CustomerRepository
	}
)

func (handler CreateOrChangeCustomerHandler) Handle(ctx context.Context, req CreateOrChangeCustomerRequest) (customer.State, error) {
	state := req.State
	if state.Id != uuid.Nil {
		var err error
		state, err = handler.repository.GetById(ctx, state.Id)
		if err != nil {
			return state, err
		}
	}

	if len(state.Name) == 0 {
		return state, ErrCustomerNameIsEmpty
	}

	if len(state.Phones) == 0 {
		return state, ErrCustomerPhonesIsEmpty
	}

	for _, phone := range state.Phones {
		size := 0
		for index, c := range phone.Number {
			if (index == 0 && c == '+') || c == ' ' {
				continue
			}

			if unicode.IsDigit(c) {
				size++
			}

			return state, ErrCustomerPhoneIsInvalid
		}

		if size < 8 || size > 12 {
			return state, ErrCustomerPhoneIsInvalid
		}
	}

	return handler.repository.Save(ctx, state)
}
