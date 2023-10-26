package products

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"happyday/infra"
	"testing"
)

func TestCreateOrChangeShouldReturnErrProductNotFound(t *testing.T) {
	id := uuid.New()

	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, id).
		Return(Product{}, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{ID: id})

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotExists, err)
}

func TestCreateOrChangeShouldReturnErrNameIsEmpty(t *testing.T) {
	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:   uuid.Nil,
		Name: "",
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsEmpty, err)
}

func TestCreateOrChangeShouldReturnErrNameTooLarge(t *testing.T) {
	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:   uuid.Nil,
		Name: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in",
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrNameTooLarge, err)
}

func TestCreateOrChangeShouldReturnErrPriceIsInvalid(t *testing.T) {
	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:    uuid.Nil,
		Name:  "Lorem ipsum",
		Price: -1,
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrPriceIsInvalid, err)
}

func TestCreateOrChangeShouldReturnErrBoxProductNotExists(t *testing.T) {
	boxID := uuid.New()
	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	repository.
		On("Exists", mock.Anything, boxID).
		Return(false, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:    uuid.Nil,
		Name:  "Lorem ipsum",
		Price: 10,
		Products: []BoxProduct{
			{ID: boxID, Quantity: 10},
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrBoxProductNotExists, err)
}

func TestCreateOrChange(t *testing.T) {
	repository := new(MockProductRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	repository.
		On("Exists", mock.Anything, mock.Anything).
		Return(true, nil)

	repository.
		On("Save", mock.Anything, mock.Anything).
		Return(Product{}, nil)

	command := &Command{repository: repository}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:    uuid.Nil,
		Name:  "Lorem ipsum",
		Price: 10,
		Products: []BoxProduct{
			{ID: uuid.New(), Quantity: 10},
		},
	})

	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	repository := new(MockProductRepository)
	repository.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	command := &Command{repository: repository}
	err := command.Delete(context.Background(), uuid.New())

	assert.Nil(t, err)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter ProductFilter) (infra.Page[Product], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(infra.Page[Product]), args.Error(1)
}

func (m *MockProductRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *MockProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductRepository) Save(ctx context.Context, product Product) (Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(Product), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
