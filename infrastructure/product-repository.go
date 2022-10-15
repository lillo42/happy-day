package infrastructure

import (
	"context"
	"errors"
	"happy_day/common"
	"happy_day/domain/product"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ProductCollection = "products"

	ProductIdAsc     ProductSortBy = "id_asc"
	ProductIdDesc    ProductSortBy = "id_desc"
	ProductNameAsc   ProductSortBy = "name_asc"
	ProductNameDesc  ProductSortBy = "name_desc"
	ProductPriceAsc  ProductSortBy = "price_asc"
	ProductPriceDesc ProductSortBy = "price_desc"
)

var (
	_ ProductRepository = (*MockProductRepository)(nil)
	_ ProductRepository = (*MongoDbProductRepository)(nil)

	ErrProductConcurrencyIssue = errors.New("product concurrency issue")
	ErrProductNotFound         = errors.New("product not found")
	ErrOneProductNotFound      = errors.New("one product in the list not found")
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

func (repository MongoDbProductRepository) GetByProducts(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
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
	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	if len(productsId) != len(products) {
		return nil, ErrOneProductNotFound
	}

	return products, nil
}

func (repository MongoDbProductRepository) GetComposed(ctx context.Context, productsId []uuid.UUID) ([]product.State, error) {
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
		err = cursor.Decode(&state)
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

func (repository MongoDbProductRepository) GetById(ctx context.Context, id uuid.UUID) (product.State, error) {
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
	if err == mongo.ErrNoDocuments {
		return product.State{}, ErrProductNotFound
	}
	if err != nil {
		return product.State{}, err
	}

	var state product.State
	err = decode.Decode(&state)
	return state, nil
}

func (repository MongoDbProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
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

func (repository MongoDbProductRepository) ExistAnyWithProduct(ctx context.Context, productId uuid.UUID) (bool, error) {
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

func (repository MongoDbProductRepository) Save(ctx context.Context, state product.State) (product.State, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return product.State{}, err
	}

	defer client.Disconnect(ctx)
	collection := client.Database(Database).
		Collection(ProductCollection)

	if state.Id == uuid.Nil {
		state.Id = uuid.New()
		state.CreateAt = time.Now().UTC()
		state.ModifyAt = time.Now().UTC()
		_, err = collection.InsertOne(ctx, state)
		return state, err
	}

	lastChange := state.ModifyAt
	state.ModifyAt = time.Now().UTC()

	res, err := collection.ReplaceOne(ctx, bson.M{"id": state.Id, "modifyAt": lastChange}, state)
	if res.ModifiedCount == 0 {
		return state, ErrProductConcurrencyIssue
	}

	return state, err
}

func (repository MongoDbProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
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
	opt := options.Find().
		SetSkip((filter.Page - 1) * filter.Size).
		SetLimit(filter.Size)

	if filter.SortBy == ProductIdAsc {
		opt.SetSort(bson.D{{"id", 1}})
	} else if filter.SortBy == ProductIdDesc {
		opt.SetSort(bson.D{{"id", -1}})
	} else if filter.SortBy == ProductNameAsc {
		opt.SetSort(bson.D{{"name", 1}})
	} else if filter.SortBy == ProductNameDesc {
		opt.SetSort(bson.D{{"name", -1}})
	} else if filter.SortBy == ProductPriceAsc {
		opt.SetSort(bson.D{{"price", 1}})
	} else if filter.SortBy == ProductPriceDesc {
		opt.SetSort(bson.D{{"price", -1}})
	}

	query := bson.M{}
	if len(filter.Text) > 0 {
		query["$or"] = []interface{}{
			bson.M{"id": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"name": bson.M{"$regex": filter.Text, "$options": "im"}},
		}
	}

	client, err := repository.CreateClient(ctx)
	var page Page[product.State]
	if err != nil {
		return page, err
	}

	collection := client.
		Database(Database).
		Collection(ProductCollection)

	totalElements, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return page, err
	}

	page.TotalElements = totalElements
	if totalElements > 0 {
		tmp := float64(totalElements) / float64(filter.Size)
		tmp = math.Ceil(tmp)
		page.TotalPages = int64(tmp)
	}

	cursor, err := collection.Find(ctx, query, opt)
	if err != nil {
		return page, err
	}

	err = cursor.All(ctx, &page.Items)
	return page, err
}
