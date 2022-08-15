package test

import (
	"happyday/abstract"
	"happyday/customer/domain"

	"github.com/stretchr/testify/mock"
)

type MockAggregateRoot struct {
	mock.Mock
}

var _ domain.AggregateRoot = (*MockAggregateRoot)(nil)

func (mock *MockAggregateRoot) State() *domain.State {
	args := mock.Called()
	return args.Get(0).(*domain.State)
}

func (mock *MockAggregateRoot) Version() int64 {
	args := mock.Called()
	return args.Get(0).(int64)
}

func (mock *MockAggregateRoot) Events() []abstract.Event {
	args := mock.Called()
	return args.Get(0).([]abstract.Event)
}

func (mock *MockAggregateRoot) Create(name, comment string, phones []domain.Phone) error {
	args := mock.Called(name, comment, phones)
	return args.Error(0)
}

func (mock *MockAggregateRoot) ChangeName(name string) error {
	args := mock.Called(name)
	return args.Error(0)
}

func (mock *MockAggregateRoot) ChangeComment(comment string) error {
	args := mock.Called(comment)
	return args.Error(0)
}

func (mock *MockAggregateRoot) ChangePhones(phones []domain.Phone) error {
	args := mock.Called(phones)
	return args.Error(0)
}
