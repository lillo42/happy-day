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
	err := validateCustomer(req.State)
	if err != nil {
		return req.State, err
	}

	if req.Id != uuid.Nil {
		state, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return state, err
		}

		req.CreatedAt = state.CreatedAt
		req.ModifiedAt = state.ModifiedAt
	}

	return handler.repository.Save(ctx, req.State)
}

func validateCustomer(state customer.State) error {
	if len(state.Name) == 0 {
		return ErrCustomerNameIsEmpty
	}

	if len(state.Phones) == 0 {
		return ErrCustomerPhonesIsEmpty
	}

	for _, phone := range state.Phones {
		size := 0
		for index, c := range phone.Number {
			if (index == 0 && c == '+') || c == ' ' {
				continue
			}

			if unicode.IsDigit(c) {
				size++
				continue
			}

			return ErrCustomerPhoneIsInvalid
		}

		if size < 8 || size > 12 {
			return ErrCustomerPhoneIsInvalid
		}
	}

	return nil
}
