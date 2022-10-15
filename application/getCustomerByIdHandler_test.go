package application

import (
	"context"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCustomerById(t *testing.T) {
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, nil)

	handler := GetCustomerByIdHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), uuid.New())
	assert.Nil(t, err)
}
