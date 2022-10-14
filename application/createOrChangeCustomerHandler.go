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
	if len(req.Name) == 0 {
		return customer.State{}, ErrCustomerNameIsEmpty
	}

	if len(req.Phones) == 0 {
		return customer.State{}, ErrCustomerPhonesIsEmpty
	}

	for _, phone := range req.Phones {
		size := 0
		for index, c := range phone.Number {
			if (index == 0 && c == '+') || c == ' ' {
				continue
			}

			if unicode.IsDigit(c) {
				size++
				continue
			}

			return customer.State{}, ErrCustomerPhoneIsInvalid
		}

		if size < 8 || size > 12 {
			return customer.State{}, ErrCustomerPhoneIsInvalid
		}
	}

	if req.Id != uuid.Nil {
		state, err := handler.repository.GetById(ctx, req.Id)
		if err != nil {
			return state, err
		}

		req.CreateAt = state.CreateAt
		req.ModifyAt = state.ModifyAt
	}

	return handler.repository.Save(ctx, req.State)
}
