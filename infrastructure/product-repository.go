package infrastructure

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ProductCollection = "products"

	IdAsc     ProductSortBy = "id_asc"
	IdDesc    ProductSortBy = "id_desc"
	NameAsc   ProductSortBy = "name_asc"
	NameDesc  ProductSortBy = "name_desc"
	PriceAsc  ProductSortBy = "price_asc"
	PriceDesc ProductSortBy = "price_desc"
)

var (
	_ ProductRepository = (*MockProductRepository)(nil)
	_ ProductRepository = (*MongoDbProductRepository)(nil)

	ErrOneProductNotFound = errors.New("one product in the list not found")
)

type (
	ProductSortBy string
	ProductFilter struct {
		Text   string
		Page   int64
		Size   int64
		SortBy ProductSortBy
	}

	ProductRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (product.State, error)
		GetByProducts(ctx context.Context, productsId []uuid.UUID) ([]product.State, error)
		GetComposed(ctx context.Context, productsId []uuid.UUID) ([]product.State, error)
		GetAll(ctx context.Context, filter ProductFilter) (Page[product.State], error)

		Exists(ctx context.Context, id uuid.UUID) (bool, error)
		ExistAnyWithProduct(ctx context.Context, productId uuid.UUID) (bool, error)
		Save(ctx context.Context, state product.State) (product.State, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	MongoDbProductRepository struct {
		MongoDbRepository
	}

	MockProductRepository struct {
		mock.Mock
	}
)

func (m *MockProductRepository) GetByProducts(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
	args := m.Called(ctx, productsId)
	return args.Get(0).([]product.State), args.Error(1)
}

func (m *MockProductRepository) GetComposed(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
	args := m.Called(ctx, productsId)
	return args.Get(0).([]product.State), args.Error(1)
}

func (m *MockProductRepository) GetById(ctx context.Context, id uuid.UUID) (product.State, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(product.State), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter ProductFilter) (Page[product.State], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(Page[product.State]), args.Error(1)
}

func (m *MockProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductRepository) ExistAnyWithProduct(ctx context.Context, productId uuid.UUID) (bool, error) {
	args := m.Called(ctx, productId)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductRepository) Save(ctx context.Context, state product.State) (product.State, error) {
	args := m.Called(ctx, state)
	return args.Get(0).(product.State), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (repository MongoDbRepository) GetByProducts(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
	query := bson.M{"id": bson.M{"$in": productsId}}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)
	cursor, err := client.Database(Database).
		Collection(ProductCollection).
		Find(ctx, query)

	if err != nil {
		return nil, err
	}

	products := make([]product.State, 0)
	for cursor.Next(ctx) {
		var state product.State
		err := cursor.Decode(&state)
		if err != nil {
			return nil, err
		}

		products = append(products, state)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	if len(productsId) != len(products) {
		return nil, ErrOneProductNotFound
	}

	return products, nil
}

func (repository MongoDbRepository) GetComposed(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
	query := bson.M{"products.id": bson.M{"$in": productsId}}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)
	cursor, err := client.Database(Database).
		Collection(ProductCollection).
		Find(ctx, query)

	if err != nil {
		return nil, err
	}

	products := make([]product.State, 0)
	for cursor.Next(ctx) {
		var state product.State
		err := cursor.Decode(&state)
		if err != nil {
			return nil, err
		}

		products = append(products, state)
	}
	if err != nil {
		return nil, err
	}

	return common.Filter(products, func(item product.State) bool {
		for _, current := range item.Products {
			found := false
			for _, expected := range productsId {
				if expected == current.Id {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	}), nil
}

func (repository MongoDbRepository) GetById(ctx context.Context, id uuid.UUID) (product.State, error) {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return product.State{}, err
	}

	defer client.Disconnect(ctx)
	decode := client.Database(Database).
		Collection(ProductCollection).
		FindOne(ctx, query)

	err = decode.Err()
	if err != nil {
		return product.State{}, err
	}

	var state product.State
	err = decode.Decode(&state)
	return state, nil
}

func (repository MongoDbRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return false, err
	}

	defer client.Disconnect(ctx)
	count, err := client.Database(Database).
		Collection(ProductCollection).
		CountDocuments(ctx, query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository MongoDbRepository) ExistAnyWithProduct(ctx context.Context, productId uuid.UUID) (bool, error) {
	query := bson.M{"products.id": productId}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return false, err
	}

	defer client.Disconnect(ctx)
	count, err := client.Database(Database).
		Collection(ProductCollection).
		CountDocuments(ctx, query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository MongoDbRepository) Save(ctx context.Context, state product.State) (product.State, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return product.State{}, err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(Database).
		Collection(ProductCollection).
		InsertOne(ctx, state)

	if err != nil {
		return product.State{}, err
	}

	return state, nil
}

func (repository MongoDbRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(Database).
		Collection(ProductCollection).
		DeleteOne(ctx, query)
	return err

}

func (repository MongoDbProductRepository) GetAll(ctx context.Context, filter ProductFilter) (Page[product.State], error) {
	return Page[product.State]{}, nil
}
