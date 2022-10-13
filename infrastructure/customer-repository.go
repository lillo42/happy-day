package infrastructure

import (
	"context"
	"happy_day/domain/customer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CustomersCollection = "customers"

	CustomerIdAsc    CustomerSortBy = "id_asc"
	CustomerIdDesc   CustomerSortBy = "id_desc"
	CustomerNameAsc  CustomerSortBy = "name_asc"
	CustomerNameDesc CustomerSortBy = "name_desc"
)

var (
	_ CustomerRepository = (*MockCustomerRepository)(nil)
	_ CustomerRepository = (*MongoDbCustomerRepository)(nil)
)

type (
	CustomerSortBy string
	CustomerFilter struct {
		Text   string
		Page   int64
		Size   int64
		SortBy CustomerSortBy
	}

	CustomerRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (customer.State, error)
		GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error)

		Save(ctx context.Context, state customer.State) (customer.State, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	MongoDbCustomerRepository struct {
		MongoDbRepository
	}

	MockCustomerRepository struct {
		mock.Mock
	}
)

func (repository *MockCustomerRepository) GetById(ctx context.Context, id uuid.UUID) (customer.State, error) {
	args := repository.Called(ctx, id)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error) {
	args := repository.Called(ctx, filter)
	return args.Get(0).(Page[customer.State]), args.Error(1)
}

func (repository *MockCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	args := repository.Called(ctx, state)
	return args.Get(0).(customer.State), args.Error(1)
}

func (repository *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := repository.Called(ctx, id)
	return args.Error(0)
}

func (repository MongoDbCustomerRepository) GetAll(ctx context.Context, filter CustomerFilter) (Page[customer.State], error) {
	return Page[customer.State]{}, nil
}

func (repository MongoDbCustomerRepository) GetById(ctx context.Context, id uuid.UUID) (customer.State, error) {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return customer.State{}, err
	}

	defer client.Disconnect(ctx)
	decode := client.Database(Database).
		Collection(CustomersCollection).
		FindOne(ctx, query)

	err = decode.Err()
	if err != nil {
		return customer.State{}, err
	}

	var state customer.State
	err = decode.Decode(&state)
	return state, nil
}

func (repository MongoDbCustomerRepository) Save(ctx context.Context, state customer.State) (customer.State, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return customer.State{}, err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(Database).
		Collection(CustomersCollection).
		InsertOne(ctx, state)

	if err != nil {
		return customer.State{}, err
	}

	return state, nil
}

func (repository MongoDbCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(Database).
		Collection(CustomersCollection).
		DeleteOne(ctx, query)
	return err

}
