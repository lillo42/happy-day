package test

import (
	"happyday/abstract"
	"happyday/product/domain"

	"github.com/stretchr/testify/mock"
)

var _ domain.AggregateRoot = (*MockAggregateRoot)(nil)

type MockAggregateRoot struct {
	mock.Mock
}

func (repository *MockAggregateRoot) State() *domain.State {
	args := repository.Called()
	return args.Get(0).(*domain.State)
}

func (repository *MockAggregateRoot) Version() int64 {
	args := repository.Called()
	return args.Get(0).(int64)
}

func (repository *MockAggregateRoot) Events() []abstract.Event {
	args := repository.Called()
	return args.Get(0).([]abstract.Event)
}

func (repository *MockAggregateRoot) Create(name string, price float64, priority int64, products []domain.Product, checkIfExists func(domain.Product) (bool, error)) error {
	args := repository.Called(name, price, priority, products, checkIfExists)
	return args.Error(0)
}

func (repository *MockAggregateRoot) ChangeName(name string) error {
	args := repository.Called(name)
	return args.Error(0)
}

func (repository *MockAggregateRoot) ChangePrice(price float64) error {
	args := repository.Called(price)
	return args.Error(0)
}

func (repository *MockAggregateRoot) Enable() error {
	args := repository.Called()
	return args.Error(0)
}

func (repository *MockAggregateRoot) Disable() error {
	args := repository.Called()
	return args.Error(0)
}

func (repository *MockAggregateRoot) ChangePriority(priority int64) error {
	args := repository.Called(priority)
	return args.Error(0)
}

func (repository *MockAggregateRoot) ChangeProducts(products []domain.Product, checkIfExists func(domain.Product) (bool, error)) error {
	args := repository.Called(products, checkIfExists)
	return args.Error(0)
}
