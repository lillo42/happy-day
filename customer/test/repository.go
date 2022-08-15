package test

import (
	"context"

	"happyday/common"
	"happyday/customer/domain"
	"happyday/customer/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

var _ infrastructure.Repository = (*MockRepository)(nil)
var _ infrastructure.ReadOnlyRepository = (*MockRepository)(nil)

func (mock *MockRepository) GetAll(ctx context.Context, filter infrastructure.Filter) (common.Page[infrastructure.ViewModel], error) {
	args := mock.Called(ctx, filter)
	return args.Get(0).(common.Page[infrastructure.ViewModel]), args.Error(1)
}

func (mock *MockRepository) GetById(ctx context.Context, id uuid.UUID) (infrastructure.DetailsViewModel, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).(infrastructure.DetailsViewModel), args.Error(1)
}

func (mock *MockRepository) Create(id uuid.UUID) domain.AggregateRoot {
	args := mock.Called(id)
	return args.Get(0).(domain.AggregateRoot)
}

func (mock *MockRepository) Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).(domain.AggregateRoot), args.Error(1)
}

func (mock *MockRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	args := mock.Called(ctx, root)
	return args.Error(0)
}

func (mock *MockRepository) Delete(ctx context.Context, root domain.AggregateRoot) error {
	args := mock.Called(ctx, root)
	return args.Error(0)
}
