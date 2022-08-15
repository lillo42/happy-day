package applications

import (
	"context"
	"errors"
	"testing"

	"happyday/common"
	"happyday/customer/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDelete_Should_ReturnError_When_ErrorToGet(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	expectingError := errors.New(common.RandString(10))
	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, expectingError)

	operation := DeleteOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestDelete_Should_ReturnError_When_ErrorToDelete(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	expectingError := errors.New(common.RandString(10))
	repository.
		On("Delete", context.Background(), root).
		Return(expectingError)

	operation := DeleteOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestDelete(t *testing.T) {
	req := DeleteRequest{
		ID: uuid.New(),
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	repository.
		On("Delete", context.Background(), root).
		Return(nil)

	operation := DeleteOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
}
