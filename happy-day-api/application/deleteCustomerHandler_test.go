package application

import (
	"happy_day/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCustomerHandler(t *testing.T) {
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := DeleteCustomerHandler{repository: repo}
	err := handler.Handle(nil, DeleteCustomerRequest{})
	assert.Nil(t, err)
}
