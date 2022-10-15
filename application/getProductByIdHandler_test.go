package application

import (
	"context"
	"happy_day/domain/product"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductById(t *testing.T) {
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(product.State{}, nil)

	handler := GetProductByIdHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), uuid.New())
	assert.Nil(t, err)
}
