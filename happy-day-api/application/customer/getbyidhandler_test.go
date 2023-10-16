package customer

import (
	"context"
	"testing"

	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCustomerById(t *testing.T) {
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, nil)

	handler := GetByIdHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), 1)
	assert.Nil(t, err)
}
