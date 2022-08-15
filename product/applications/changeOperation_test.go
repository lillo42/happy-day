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

func TestChange_Should_ReturnError_When_ErrorToGet(t *testing.T) {
	req := ChangeRequest{
		ID: uuid.New(),
	}

	expectedError := errors.New(common.RandString(10))
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return((*test.MockAggregateRoot)(nil), expectedError)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangeName(t *testing.T) {
	name := common.RandString(10)
	req := ChangeRequest{
		ID:   uuid.New(),
		Name: name,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("ChangeName", name).
		Return(expectedError)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangePrice(t *testing.T) {
	price := 10.0
	req := ChangeRequest{
		ID:    uuid.New(),
		Price: price,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", price).
		Return(expectedError)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangePriority(t *testing.T) {
	priority := int64(1)
	req := ChangeRequest{
		ID:       uuid.New(),
		Priority: priority,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", priority).
		Return(expectedError)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangeProducts(t *testing.T) {
	var products []domain.Product
	req := ChangeRequest{
		ID:       uuid.New(),
		Products: products,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangeProducts", products, mock.Anything).
		Return(expectedError)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToEnable(t *testing.T) {
	req := ChangeRequest{
		ID:       uuid.New(),
		IsEnable: true,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Enable").
		Return(expectedError)

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangeProducts", mock.Anything, mock.Anything).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToDisable(t *testing.T) {
	req := ChangeRequest{
		ID:       uuid.New(),
		IsEnable: false,
	}

	expectedError := errors.New(common.RandString(10))
	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Disable").
		Return(expectedError)

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangeProducts", mock.Anything, mock.Anything).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange_Should_ReturnError_When_ErrorToSave(t *testing.T) {
	req := ChangeRequest{
		ID: uuid.New(),
	}

	aggregateRoot := &test.MockAggregateRoot{}
	aggregateRoot.
		On("Disable").
		Return(nil)

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangeProducts", mock.Anything, mock.Anything).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	expectedError := errors.New(common.RandString(10))
	repository.
		On("Save", context.Background(), aggregateRoot).
		Return(expectedError)

	operation := &ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Equal(t, expectedError, err)
}

func TestChange(t *testing.T) {
	req := ChangeRequest{
		ID:       uuid.New(),
		IsEnable: true,
	}

	aggregateRoot := &test.MockAggregateRoot{}

	aggregateRoot.
		On("Enable").
		Return(nil)

	aggregateRoot.
		On("ChangeName", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePrice", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangePriority", mock.Anything).
		Return(nil)

	aggregateRoot.
		On("ChangeProducts", mock.Anything, mock.Anything).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(aggregateRoot, nil)

	repository.
		On("Save", context.Background(), aggregateRoot).
		Return(nil)

	operation := &ChangeOperation{
		repository: repository,
	}

	res, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, common.VoidResponse{}, res)
}
