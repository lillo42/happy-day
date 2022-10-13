package infrastructure

import (
	"context"
	"happy_day/domain/customer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

const (
	CustomerIdAsc    CustomerSortBy = "id_asc"
	CustomerIdDesc   CustomerSortBy = "id_desc"
	CustomerNameAsc  CustomerSortBy = "name_asc"
	CustomerNameDesc CustomerSortBy = "name_desc"
)

var (
	_ CustomerRepository = (*MockCustomerRepository)(nil)
)

type (
	CustomerSortBy string
	CustomerFilter struct {
		Text   string
		Page   int64
		Size   int64
		SortBy CustomerSortBy
	}
	CustomerRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (customer.State, error)
		GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error)

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

func (repository *MockCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error) {
	args := repository.Called(ctx, filter)
	return args.Get(0).(Page[customer.State]), args.Error(1)
}

func (repository *MockCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	args := repository.Called(ctx, state)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := repository.Called(ctx, id)
	return args.Error(0)
}
