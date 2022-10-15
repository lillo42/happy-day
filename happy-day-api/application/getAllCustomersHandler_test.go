package application

import (
	"context"
	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllCustomers(t *testing.T) {
	req := infrastructure.CustomerFilter{
		Page: 0,
		Size: 0,
		Text: common.RandString(10),
	}

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[customer.State]{}, nil)

	handler := GetAllCustomersHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
