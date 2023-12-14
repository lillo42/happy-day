package customers

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mec/infra"
	"testing"
)

func TestCreateOrChangeNameIsEmpty(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsEmpty, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangeNameIsLargeThan255(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		Name: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in",
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsTooLarge, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangePhoneIsEmpty(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		Name:   "Lorem ipsum dolor",
		Phones: []string{""},
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrPhoneNumberIsInvalid, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangePhoneIsNotNumber(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		Name:   "Lorem ipsum dolor",
		Phones: []string{"abcdefgh"},
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrPhoneNumberIsInvalid, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangePhoneIsTooShort(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		Name:   "Lorem ipsum dolor",
		Phones: []string{"1234567"},
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrPhoneNumberIsInvalid, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangePhoneIsTooLarge(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		Name:   "Lorem ipsum dolor",
		Phones: []string{"123456789012"},
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrPhoneNumberIsInvalid, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChangeNotFound(t *testing.T) {
	repository := new(MockCustomerRepository)
	repository.On("GetOrCreate", mock.Anything, mock.Anything).Return(Customer{
		ID: uuid.New(),
	}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		ID:   uuid.New(),
		Name: "Lorem ipsum dolor",
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err)

	repository.AssertExpectations(t)
}

func TestCreateOrChange(t *testing.T) {
	id := uuid.New()

	repository := new(MockCustomerRepository)
	repository.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Customer{ID: id}, nil)

	repository.
		On("Save", mock.Anything, Customer{ID: id, Name: "Lorem ipsum dolor"}).
		Return(Customer{ID: id}, nil)

	command := &Command{
		repository: repository,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChangeCustomer{
		ID:   id,
		Name: "Lorem ipsum dolor",
	})

	assert.Nil(t, err)

	repository.AssertExpectations(t)
}

var _ CustomerRepository = (*MockCustomerRepository)(nil)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (infra.Page[Customer], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(infra.Page[Customer]), args.Error(1)
}

func (m *MockCustomerRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Customer), args.Error(1)
}

func (m *MockCustomerRepository) Save(ctx context.Context, customer Customer) (Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(Customer), args.Error(1)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
