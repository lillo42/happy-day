package applications

import (
	"context"
	"errors"
	"testing"

	"happyday/common"
	"happyday/product/domain"
	"happyday/product/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Should_ReturnError_When_CreateReturnError(t *testing.T) {
	req := CreateRequest{
		Name:     common.RandString(10),
		Price:    100.0,
		Priority: 1,
		Products: []domain.Product{},
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Create", req.Name, req.Price, req.Priority, req.Products, mock.Anything).
		Return(expectedError)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(aggregateRoot)

	operation := CreateOperation{repository: repository}
	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, err, expectedError)
}

func TestCreate_Should_ReturnError_When_ErrorToSave(t *testing.T) {
	req := CreateRequest{
		Name:     common.RandString(10),
		Price:    100.0,
		Priority: 1,
		Products: []domain.Product{},
	}

	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Create", req.Name, req.Price, req.Priority, req.Products, mock.Anything).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(aggregateRoot)

	expectedError := errors.New(common.RandString(10))
	repository.
		On("Save", context.Background(), aggregateRoot).
		Return(expectedError)

	operation := CreateOperation{repository: repository}
	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, err, expectedError)
}

func TestCreate(t *testing.T) {
	req := CreateRequest{
		Name:     common.RandString(10),
		Price:    100.0,
		Priority: 1,
		Products: []domain.Product{},
	}

	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Create", req.Name, req.Price, req.Priority, req.Products, mock.Anything).
		Return(nil)

	state := domain.NewStateWithID(uuid.New())
	aggregateRoot.
		On("State").
		Return(state)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(aggregateRoot)

	repository.
		On("Save", context.Background(), aggregateRoot).
		Return(nil)

	operation := CreateOperation{repository: repository}
	res, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, res.ID, state.ID())
}
