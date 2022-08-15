package applications

import (
	"context"
	"errors"
	"testing"

	"happyday/common"
	"happyday/customer/domain"
	"happyday/customer/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Should_ReturnError_When_ErrorToCreate(t *testing.T) {
	req := CreateRequest{
		Name:    common.RandString(10),
		Comment: "",
		Phones:  []domain.Phone{},
	}

	expectingError := errors.New(common.RandString(10))
	root := &test.MockAggregateRoot{}
	root.
		On("Create", req.Name, req.Comment, req.Phones).
		Return(expectingError)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(root)

	operation := CreateOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestCreate_Should_ReturnError_When_ErrorToSave(t *testing.T) {
	req := CreateRequest{
		Name:    common.RandString(10),
		Comment: "",
		Phones:  []domain.Phone{},
	}

	root := &test.MockAggregateRoot{}
	root.
		On("Create", req.Name, req.Comment, req.Phones).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(root)

	expectingError := errors.New(common.RandString(10))
	repository.
		On("Save", context.Background(), root).
		Return(expectingError)

	operation := CreateOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestCreate(t *testing.T) {
	req := CreateRequest{
		Name:    common.RandString(10),
		Comment: "",
		Phones:  []domain.Phone{},
	}

	root := &test.MockAggregateRoot{}
	root.
		On("Create", req.Name, req.Comment, req.Phones).
		Return(nil)

	repository := &test.MockRepository{}
	repository.
		On("Create", mock.Anything).
		Return(root)

	repository.
		On("Save", context.Background(), root).
		Return(nil)

	operation := CreateOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
}
