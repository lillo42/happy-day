package applications

import (
	"context"
	"errors"
	"testing"

	"happyday/common"
	"happyday/customer/domain"
	"happyday/customer/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChange_Should_ReturnError_When_ErrorToGet(t *testing.T) {
	req := ChangeRequest{
		ID: uuid.New(),
	}

	expectingError := errors.New(common.RandString(10))
	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, expectingError)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangeName(t *testing.T) {
	req := ChangeRequest{
		ID:   uuid.New(),
		Name: common.RandString(10),
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	expectingError := errors.New(common.RandString(10))
	root.
		On("ChangeName", req.Name).
		Return(expectingError)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangeComment(t *testing.T) {
	req := ChangeRequest{
		ID:      uuid.New(),
		Name:    common.RandString(10),
		Comment: common.RandString(10),
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	expectingError := errors.New(common.RandString(10))
	root.
		On("ChangeName", req.Name).
		Return(nil)

	root.
		On("ChangeComment", req.Comment).
		Return(expectingError)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestChange_Should_ReturnError_When_ErrorToChangePhones(t *testing.T) {
	req := ChangeRequest{
		ID:      uuid.New(),
		Name:    common.RandString(10),
		Comment: common.RandString(10),
		Phones: []domain.Phone{
			domain.NewPhone(common.RandString(10)),
		},
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	expectingError := errors.New(common.RandString(10))
	root.
		On("ChangeName", req.Name).
		Return(nil)

	root.
		On("ChangeComment", req.Comment).
		Return(nil)

	root.
		On("ChangePhones", req.Phones).
		Return(expectingError)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestChange_Should_ReturnError_When_ErrorToSave(t *testing.T) {
	req := ChangeRequest{
		ID:      uuid.New(),
		Name:    common.RandString(10),
		Comment: common.RandString(10),
		Phones: []domain.Phone{
			domain.NewPhone(common.RandString(10)),
		},
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	root.
		On("ChangeName", req.Name).
		Return(nil)

	root.
		On("ChangeComment", req.Comment).
		Return(nil)

	root.
		On("ChangePhones", req.Phones).
		Return(nil)

	expectingError := errors.New(common.RandString(10))
	repository.
		On("Save", context.Background(), root).
		Return(expectingError)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectingError, err)
}

func TestChange(t *testing.T) {
	req := ChangeRequest{
		ID:      uuid.New(),
		Name:    common.RandString(10),
		Comment: common.RandString(10),
		Phones: []domain.Phone{
			domain.NewPhone(common.RandString(10)),
		},
	}

	root := &test.MockAggregateRoot{}
	repository := &test.MockRepository{}
	repository.
		On("Get", context.Background(), req.ID).
		Return(root, nil)

	root.
		On("ChangeName", req.Name).
		Return(nil)

	root.
		On("ChangeComment", req.Comment).
		Return(nil)

	root.
		On("ChangePhones", req.Phones).
		Return(nil)

	repository.
		On("Save", context.Background(), root).
		Return(nil)

	operation := ChangeOperation{
		repository: repository,
	}

	_, err := operation.Execute(context.Background(), req)
	assert.Nil(t, err)
}
