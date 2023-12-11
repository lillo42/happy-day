package customers

import (
	"context"
	"github.com/google/uuid"
	"strconv"
)

type (
	Command struct {
		repository CustomerRepository
	}

	CreateOrChangeCustomer struct {
		ID      uuid.UUID `json:"id,omitempty"`
		Name    string    `json:"name"`
		Comment string    `json:"comment,omitempty"`
		Phones  []string  `json:"phones,omitempty"`
	}
)

func (command *Command) CreateOrChange(ctx context.Context, req CreateOrChangeCustomer) (Customer, error) {
	customer, err := command.repository.GetOrCreate(ctx, req.ID)

	if err != nil {
		return Customer{}, err
	}

	if req.ID != uuid.Nil && req.ID != customer.ID {
		return Customer{}, ErrNotFound
	}

	if len(req.Name) == 0 {
		return Customer{}, ErrNameIsEmpty
	}

	if len(req.Name) > 255 {
		return Customer{}, ErrNameIsTooLarge
	}

	for _, phone := range req.Phones {
		if len(phone) < 8 || len(phone) > 11 {
			return Customer{}, ErrPhoneNumberIsInvalid
		}

		_, err := strconv.ParseUint(phone, 10, 64)
		if err != nil {
			return Customer{}, ErrPhoneNumberIsInvalid
		}
	}

	customer.Name = req.Name
	customer.Comment = req.Comment
	customer.Phones = req.Phones

	return command.repository.Save(ctx, customer)
}

func (command *Command) Delete(ctx context.Context, id uuid.UUID) error {
	return command.repository.Delete(ctx, id)
}
