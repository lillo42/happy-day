package discounts

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mec/infra"
	"testing"
)

func TestCreateOrChangeWhenNameIsInvalid(t *testing.T) {
	repo := new(MockDiscountRepository)

	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	service := new(MockProductService)

	command := &Command{
		repository:     repo,
		productService: service,
	}

	t.Run("name is empty", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID: uuid.Nil,
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrNameIsEmpty, err)
	})

	t.Run("name is too large", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:   uuid.Nil,
			Name: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in",
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrNameIsTooLarge, err)
	})
}

func TestCreateOrChangeWhenPriceIsInvalid(t *testing.T) {
	repo := new(MockDiscountRepository)

	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	service := new(MockProductService)

	command := &Command{
		repository:     repo,
		productService: service,
	}

	t.Run("price is zero", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:    uuid.Nil,
			Name:  "Lorem ipsum",
			Price: 0,
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrPriceIsInvalid, err)
	})

	t.Run("price is negative", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:    uuid.Nil,
			Name:  "Lorem ipsum",
			Price: -1,
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrPriceIsInvalid, err)
	})

}

func TestCreateOrChangeWhenProductIsInvalid(t *testing.T) {
	repo := new(MockDiscountRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	t.Run("product is empty", func(t *testing.T) {
		service := new(MockProductService)

		command := &Command{
			repository:     repo,
			productService: service,
		}

		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:       uuid.Nil,
			Name:     "Lorem ipsum",
			Price:    1,
			Products: make([]Product, 0),
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrProductsIsMissing, err)
	})

	t.Run("exists return error", func(t *testing.T) {
		expectedError := errors.New("test")
		service := new(MockProductService)
		service.
			On("Exists", mock.Anything, mock.Anything).
			Return(false, expectedError)

		command := &Command{
			repository:     repo,
			productService: service,
		}

		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:    uuid.Nil,
			Name:  "Lorem ipsum",
			Price: 1,
			Products: []Product{
				{ID: uuid.New(), Quantity: 0},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("product not exists", func(t *testing.T) {
		service := new(MockProductService)
		service.
			On("Exists", mock.Anything, mock.Anything).
			Return(false, nil)

		command := &Command{
			repository:     repo,
			productService: service,
		}

		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			ID:    uuid.Nil,
			Name:  "Lorem ipsum",
			Price: 1,
			Products: []Product{
				{ID: uuid.New(), Quantity: 0},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrProductNotFound, err)
	})
}

func TestCreateOrChangeWhenDiscountNotExists(t *testing.T) {
	repo := new(MockDiscountRepository)

	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	service := new(MockProductService)

	command := &Command{
		repository:     repo,
		productService: service,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID: uuid.New(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrDiscountNotFound, err)
}

func TestCreateOrChange(t *testing.T) {
	repo := new(MockDiscountRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(Discount{}, nil)

	service := new(MockProductService)
	service.
		On("Exists", mock.Anything, mock.Anything).
		Return(true, nil)

	command := &Command{
		repository:     repo,
		productService: service,
	}
	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID:    uuid.Nil,
		Name:  "Lorem ipsum",
		Price: 1,
		Products: []Product{
			{ID: uuid.New(), Quantity: 1},
		},
	})

	assert.Nil(t, err)
}

var (
	_ DiscountRepository = (*MockDiscountRepository)(nil)
	_ ProductService     = (*MockProductService)(nil)
)

type (
	MockDiscountRepository struct {
		mock.Mock
	}

	MockProductService struct {
		mock.Mock
	}
)

func (m *MockProductService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDiscountRepository) GetAll(ctx context.Context, filter DiscountFilter) (infra.Page[Discount], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(infra.Page[Discount]), args.Error(1)
}

func (m *MockDiscountRepository) GetAllWithProducts(ctx context.Context, productsId []uuid.UUID) ([]Discount, error) {
	args := m.Called(ctx, productsId)
	return args.Get(0).([]Discount), args.Error(1)
}

func (m *MockDiscountRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Discount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Discount), args.Error(1)
}

func (m *MockDiscountRepository) Save(ctx context.Context, discount Discount) (Discount, error) {
	args := m.Called(ctx, discount)
	return args.Get(0).(Discount), args.Error(1)
}

func (m *MockDiscountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(1)
}
