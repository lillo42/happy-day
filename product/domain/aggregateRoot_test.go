package domain

import (
	"errors"
	"testing"

	"happyday/abstract"
	"happyday/common"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChangeName_Should_ReturnError_When_NameIsEmpty(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangeName("")
	assert.NotNil(t, err)
	assert.Equal(t, NameIsEmpty, err)
}

func TestChangeName_Should_ReturnError_When_NameIsLargeThan100(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangeName(common.RandString(1000))
	assert.NotNil(t, err)
	assert.Equal(t, NameIsTooLarge, err)
}

func TestChangeName(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	name := common.RandString(10)
	err := root.ChangeName(name)
	assert.Nil(t, err)
	assert.Equal(t, name, root.State().Name())
	assert.Equal(t, 1, len(root.Events()))
}

func TestChangePrice_Should_ReturnError_When_PriceIsLessThanZero(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangePrice(-1)
	assert.NotNil(t, err)
	assert.Equal(t, PriceIsInvalid, err)
}

func TestChangePrice(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	price := 100.0
	err := root.ChangePrice(price)
	assert.Nil(t, err)
	assert.Equal(t, price, root.State().Price())
	assert.Equal(t, 1, len(root.Events()))
}

func TestEnable(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{
			isEnable: false,
		}, 0),
	}

	err := root.Enable()
	assert.Nil(t, err)
	assert.Equal(t, true, root.State().IsEnable())
	assert.Equal(t, 1, len(root.Events()))
}

func TestDisable(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{
			isEnable: true,
		}, 0),
	}

	err := root.Disable()
	assert.Nil(t, err)
	assert.Equal(t, false, root.State().IsEnable())
	assert.Equal(t, 1, len(root.Events()))
}

func TestChangePriority(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{
			priority: -1,
		}, 0),
	}

	err := root.ChangePriority(10)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), root.State().Priority())
	assert.Equal(t, 1, len(root.Events()))
}

func TestChangeProducts_Should_ReturnError_When_ErrorToExecuteFunc(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	expectedError := errors.New(common.RandString(10))
	err := root.ChangeProducts([]Product{{id: uuid.New()}}, func(product Product) (bool, error) {
		return false, expectedError
	})

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, len(root.Events()))
}

func TestChangeProducts_Should_ReturnProductNotExist_When_FuncFalse(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	err := root.ChangeProducts([]Product{{id: uuid.New()}}, func(product Product) (bool, error) {
		return false, nil
	})

	assert.NotNil(t, err)
	assert.Equal(t, ProductNotExist, err)
	assert.Equal(t, 0, len(root.Events()))
}

func TestChangeProducts(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	products := []Product{{id: uuid.New()}}
	err := root.ChangeProducts(products, func(product Product) (bool, error) { return true, nil })

	assert.Nil(t, err)
	assert.Equal(t, products, root.State().Products())
	assert.Equal(t, 1, len(root.Events()))
}

func TestCreate_Should_ReturnError_When_NameIsEmpty(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create("", 0, 0, []Product{}, func(product Product) (bool, error) { return false, nil })
	assert.NotNil(t, err)
	assert.Equal(t, NameIsEmpty, err)
}

func TestCreate_Should_ReturnError_When_NameIsLargeThan100(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(1000), 0, 0, []Product{}, func(product Product) (bool, error) { return false, nil })
	assert.NotNil(t, err)
	assert.Equal(t, NameIsTooLarge, err)
}

func TestCreate_Should_ReturnError_When_PriceIsLessThanZero(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), -1, 0, []Product{}, func(product Product) (bool, error) { return false, nil })
	assert.NotNil(t, err)
	assert.Equal(t, PriceIsInvalid, err)
}

func TestCreate_Should_ReturnError_When_ErrorToExecuteFunc(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	expectedError := errors.New(common.RandString(10))
	err := root.Create(common.RandString(10), 1, 0, []Product{{id: uuid.New()}}, func(product Product) (bool, error) {
		return false, expectedError
	})

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, len(root.Events()))
}

func TestCreate_Should_ReturnProductNotExist_When_FuncFalse(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	err := root.Create(common.RandString(10), 10, 0, []Product{{id: uuid.New()}}, func(product Product) (bool, error) {
		return false, nil
	})

	assert.NotNil(t, err)
	assert.Equal(t, ProductNotExist, err)
	assert.Equal(t, 0, len(root.Events()))
}

func TestCreate(t *testing.T) {
	root := &aggregateRoot{abstract.NewAggregateRoot(&State{}, 0)}

	name := common.RandString(10)
	price := float64(10)
	priority := int64(2)
	products := []Product{{id: uuid.New()}}
	err := root.Create(name, price, priority, products, func(product Product) (bool, error) { return true, nil })

	assert.Nil(t, err)
	assert.Equal(t, 1, len(root.Events()))
	assert.Equal(t, name, root.State().Name())
	assert.Equal(t, price, root.State().Price())
	assert.Equal(t, priority, root.State().Priority())
	assert.Equal(t, products, root.State().Products())
}
