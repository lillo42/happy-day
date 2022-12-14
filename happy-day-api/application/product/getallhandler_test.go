package product

import (
	"context"
	"testing"

	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllProducts(t *testing.T) {
	req := infrastructure.ProductFilter{
		Page: 0,
		Size: 0,
		Text: common.RandString(10),
	}

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[product.State]{}, nil)

	handler := GetAllHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
