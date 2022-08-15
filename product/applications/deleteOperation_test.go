package applications

import (
	"context"
	"errors"
	"testing"

	"happyday/common"
	"happyday/product/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOperation_Should_ReturnError_When_ErrorToGet(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	repository := &test.MockRepository{}

	expectedError := errors.New(common.RandString(10))
	repository.
		On("Get", context.Background(), req.ID).
		Return((*test.MockAggregateRoot)(nil), expectedError)

	operation := &DeleteOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestDeleteOperation_Should_ReturnError_When_ErrorToDelete(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	aggregateRoot := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}

	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	expectedError := errors.New(common.RandString(10))
	repository.
		On("Delete", context.Background(), aggregateRoot).
		Return(expectedError)

	operation := &DeleteOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestDeleteOperation(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	aggregateRoot := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}

	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	repository.
		On("Delete", context.Background(), aggregateRoot).
		Return(nil)

	operation := &DeleteOperation{
		repository: repository,
	}

	res, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, common.VoidResponse{}, res)
}
