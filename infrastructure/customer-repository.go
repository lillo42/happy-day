package infrastructure

import (
	"context"
	"happy_day/domain/customer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type (
	CustomerRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (customer.State, error)

		Save(ctx context.Context, state customer.State) (customer.State, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	MockCustomerRepository struct {
		mock.Mock
	}
)

func (repository *MockCustomerRepository) GetById(ctx context.Context, id uuid.UUID) (customer.State, error) {
	args := repository.Called(ctx, id)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	args := repository.Called(ctx, state)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := repository.Called(ctx, id)
	return args.Error(0)
}
