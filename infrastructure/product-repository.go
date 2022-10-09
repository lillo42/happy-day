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
)

var (
	_ ProductRepository = (*MockProductRepository)(nil)
	_ ProductRepository = (*MongoDbProductRepository)(nil)

	ErrOneProductNotFound = errors.New("one product in the list not found")
)

type (
	ProductRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (product.State, error)
		GetByProducts(ctx context.Context, productsId []uuid.UUID) ([]product.State, error)
		GetComposed(ctx context.Context, productsId []uuid.UUID) ([]product.State, error)
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
