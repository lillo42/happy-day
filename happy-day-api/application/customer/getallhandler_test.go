package customer

import (
	"context"
	"testing"

	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllCustomers(t *testing.T) {
	req := infrastructure.CustomerFilter{
		Page:    0,
		Size:    0,
		Name:    common.RandString(10),
		Comment: common.RandString(10),
	}

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[customer.State]{}, nil)

	handler := GetAllHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
