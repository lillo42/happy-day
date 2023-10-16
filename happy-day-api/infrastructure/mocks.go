package infrastructure

import (
	"context"
	"github.com/stretchr/testify/mock"
	"happy_day/domain/customer"
)

var (
	_ CustomerRepository = (*MockCustomerRepository)(nil)
)

type MockCustomerRepository struct {
	mock.Mock
}

func (repository *MockCustomerRepository) GetById(ctx context.Context, id uint) (customer.State, error) {
	args := repository.Called(ctx, id)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error) {
	args := repository.Called(ctx, filter)
	return args.Get(0).(Page[customer.State]), args.Error(1)
}

func (repository *MockCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	args := repository.Called(ctx, state)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) Delete(ctx context.Context, id uint) error {
	args := repository.Called(ctx, id)
	return args.Error(0)
}
