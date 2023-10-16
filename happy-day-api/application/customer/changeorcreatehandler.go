package customer

import (
	"unicode"

	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"golang.org/x/net/context"
)

type (
	ChangeOrCreateRequest struct {
		customer.State
	}

	ChangeOrCreateHandler struct {
		repository infrastructure.CustomerRepository
	}
)

func (handler ChangeOrCreateHandler) Handle(ctx context.Context, req ChangeOrCreateRequest) (customer.State, error) {
	err := Validate(req.State)
	if err != nil {
		return req.State, err
	}

	if req.ID > 0 {
		state, err := handler.repository.GetById(ctx, req.ID)
		if err != nil {
			return state, err
		}

		req.CreatedAt = state.CreatedAt
		req.UpdateAt = state.UpdateAt
	}

	return handler.repository.Save(ctx, req.State)
}

func Validate(state customer.State) error {
	if len(state.Name) == 0 {
		return infrastructure.ErrCustomerNameIsEmpty
	}

	if len(state.Phones) == 0 {
		return infrastructure.ErrCustomerPhonesIsEmpty
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

			return infrastructure.ErrCustomerPhoneIsInvalid
		}

		if size < 8 || size > 12 {
			return infrastructure.ErrCustomerPhoneIsInvalid
		}
	}

	return nil
}
