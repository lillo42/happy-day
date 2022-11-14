package product

import (
	"context"
	"errors"
	"testing"

	"happy_day/common"
	"happy_day/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteProductWhenErrCheckIfExistAnyWithProduct(t *testing.T) {
	req := DeleteRequest{
		Id: uuid.New(),
	}

	repo := &infrastructure.MockProductRepository{}
	expectedErr := errors.New(common.RandString(10))
	repo.
		On("ExistAnyWithProduct", mock.Anything, req.Id).
		Return(false, expectedErr)
	handler := DeleteHandler{
		repository: repo,
	}

	err := handler.Handle(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestDeleteProductWhenExistAnyWithProduct(t *testing.T) {
	req := DeleteRequest{
		Id: uuid.New(),
	}

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("ExistAnyWithProduct", mock.Anything, req.Id).
		Return(true, nil)

	handler := DeleteHandler{
		repository: repo,
	}

	err := handler.Handle(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrExistOtherProductWithThisProduct, err)
}

func TestDeleteProduct(t *testing.T) {
	req := DeleteRequest{
		Id: uuid.New(),
	}

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("ExistAnyWithProduct", mock.Anything, req.Id).
		Return(false, nil)

	repo.
		On("Delete", mock.Anything, req.Id).
		Return(nil)

	handler := DeleteHandler{
		repository: repo,
	}

	err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
