package application

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/infrastructure"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrChangeProductWhenProductNotFound(t *testing.T) {
	req := CreateOrChangeProductRequest{
		State: product.State{
			Id: uuid.New(),
		},
	}

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, req.Id).
		Return(product.State{}, expectedErr)

	handler := CreateOrChangeProductHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreateOrChangeProductWhenNameIsEmpty(t *testing.T) {
	req := CreateOrChangeProductRequest{
		State: product.State{
			Name: "",
		},
	}

	handler := CreateOrChangeProductHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNameIsEmpty, err)
}

func TestCreateOrChangeProductWhenPriceIsLessThanZero(t *testing.T) {
	req := CreateOrChangeProductRequest{
		State: product.State{
			Name:  common.RandString(10),
			Price: -1,
		},
	}

	handler := CreateOrChangeProductHandler{}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductPriceIsLessThanZero, err)
}

func TestCreateOrChangeProductWhenErrToCheckIfExists(t *testing.T) {
	req := CreateOrChangeProductRequest{
		State: product.State{
			Name:  common.RandString(10),
			Price: 10,
			Products: []product.Product{
				{
					Id: uuid.New(),
				},
			},
		},
	}

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("Exists", mock.Anything, mock.Anything).
		Return(false, expectedErr)

	handler := CreateOrChangeProductHandler{
		repository: repo,
	}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreateOrChangeProductWhenProductNotExsits(t *testing.T) {
	req := CreateOrChangeProductRequest{
		State: product.State{
			Name:  common.RandString(10),
			Price: 10,
			Products: []product.Product{
				{
					Id: uuid.New(),
				},
			},
		},
	}

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("Exists", mock.Anything, mock.Anything).
		Return(false, nil)

	handler := CreateOrChangeProductHandler{
		repository: repo,
	}
	_, err := handler.Handle(context.Background(), req)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestCreateOrChangeProductWhenProductAmountIsLessOrEqualThanZero(t *testing.T) {
	reqs := []CreateOrChangeProductRequest{
		{
			State: product.State{
				Name:  common.RandString(10),
				Price: 10,
				Products: []product.Product{
					{
						Id:     uuid.New(),
						Amount: -1,
					},
				},
			},
		},
		{

			State: product.State{
				Name:  common.RandString(10),
				Price: 10,
				Products: []product.Product{
					{
						Id:     uuid.New(),
						Amount: 0,
					},
				},
			},
		},
	}

	for _, req := range reqs {
		t.Run("With Amount "+strconv.FormatInt(req.Products[0].Amount, 10), func(t *testing.T) {
			repo := &infrastructure.MockProductRepository{}
			repo.
				On("Exists", mock.Anything, mock.Anything).
				Return(true, nil)

			handler := CreateOrChangeProductHandler{
				repository: repo,
			}
			_, err := handler.Handle(context.Background(), req)

			assert.NotNil(t, err)
			assert.Equal(t, ErrProductAmoutIsInvalid, err)
		})
	}
}

func TestCreateOrChangeProduct(t *testing.T) {
	reqs := []CreateOrChangeProductRequest{
		{
			State: product.State{
				Id:    uuid.New(),
				Name:  common.RandString(10),
				Price: 10,
				Products: []product.Product{
					{
						Id:     uuid.New(),
						Amount: 10,
					},
				},
			},
		},
		{

			State: product.State{
				Name:  common.RandString(10),
				Price: 10,
				Products: []product.Product{
					{
						Id:     uuid.New(),
						Amount: 20,
					},
				},
			},
		},
	}

	for _, req := range reqs {
		t.Run("With id "+req.Id.String(), func(t *testing.T) {
			repo := &infrastructure.MockProductRepository{}

			if req.Id != uuid.Nil {
				repo.
					On("GetById", mock.Anything, req.Id).
					Return(req.State, nil)
			}

			repo.
				On("Exists", mock.Anything, mock.Anything).
				Return(true, nil)
			repo.
				On("Save", mock.Anything, mock.Anything).
				Return(req.State, nil)

			handler := CreateOrChangeProductHandler{
				repository: repo,
			}
			_, err := handler.Handle(context.Background(), req)

			assert.Nil(t, err)
		})
	}
}
