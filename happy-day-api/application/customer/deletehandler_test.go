package customer

import (
	"testing"

	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCustomerHandler(t *testing.T) {
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := DeleteHandler{repository: repo}
	err := handler.Handle(nil, DeleteRequest{})
	assert.Nil(t, err)
}
