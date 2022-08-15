package test

import (
	"context"

	"happyday/common"
	"happyday/product/domain"
	"happyday/product/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var _ infrastructure.Repository = (*MockRepository)(nil)

type MockRepository struct {
	mock.Mock
}

func (repository *MockRepository) Exists(ctx context.Context, product uuid.UUID) (bool, error) {
	args := repository.Called(ctx, product)
	return args.Get(0).(bool), args.Error(1)
}

func (repository *MockRepository) GetAll(context context.Context, filter infrastructure.Filter) (common.Page[infrastructure.ViewModel], error) {
	args := repository.Called(context, filter)
	return args.Get(0).(common.Page[infrastructure.ViewModel]), args.Error(1)
}

func (repository *MockRepository) GetById(context context.Context, id uuid.UUID) (infrastructure.DetailsViewModel, error) {
	args := repository.Called(context, id)
	return args.Get(0).(infrastructure.DetailsViewModel), args.Error(1)
}

func (repository *MockRepository) Create(id uuid.UUID) domain.AggregateRoot {
	args := repository.Called(id)
	return args.Get(0).(domain.AggregateRoot)
}

func (repository *MockRepository) Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error) {
	args := repository.Called(ctx, id)
	return args.Get(0).(domain.AggregateRoot), args.Error(1)
}

func (repository *MockRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	args := repository.Called(ctx, root)
	return args.Error(0)
}

func (repository *MockRepository) Delete(ctx context.Context, root domain.AggregateRoot) error {
	args := repository.Called(ctx, root)
	return args.Error(0)
}
