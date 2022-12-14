package infrastructure

import (
	"context"
	"errors"
	"math"
	"time"

	"happy_day/domain/customer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CustomersCollection = "customers"

	CustomerIdAsc       CustomerSortBy = "id_asc"
	CustomerIdDesc      CustomerSortBy = "id_desc"
	CustomerNameAsc     CustomerSortBy = "name_asc"
	CustomerNameDesc    CustomerSortBy = "name_desc"
	CustomerCommentAsc  CustomerSortBy = "comment_asc"
	CustomerCommentDesc CustomerSortBy = "comment_desc"
)

var (
	_ CustomerRepository = (*MockCustomerRepository)(nil)
	_ CustomerRepository = (*MongoDbCustomerRepository)(nil)

	ErrCustomerConcurrencyIssue = errors.New("product concurrency issue")
	ErrCustomerNotFound         = errors.New("product not found")
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
	opt := options.Find().
		SetSkip((filter.Page - 1) * filter.Size).
		SetLimit(filter.Size)

	if filter.SortBy == CustomerIdAsc {
		opt.SetSort(bson.D{{"id", 1}})
	} else if filter.SortBy == CustomerIdDesc {
		opt.SetSort(bson.D{{"id", -1}})
	} else if filter.SortBy == CustomerNameAsc {
		opt.SetSort(bson.D{{"name", 1}})
	} else if filter.SortBy == CustomerNameDesc {
		opt.SetSort(bson.D{{"name", -1}})
	} else if filter.SortBy == CustomerCommentAsc {
		opt.SetSort(bson.D{{"comment", 1}})
	} else if filter.SortBy == CustomerCommentDesc {
		opt.SetSort(bson.D{{"comment", -1}})
	}

	query := bson.M{}
	if len(filter.Text) > 0 {
		query["$or"] = []interface{}{
			bson.M{"id": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"name": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"comment": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"phones": bson.M{"$elemMatch": bson.M{"number": bson.M{"$regex": filter.Text, "$options": "im"}}}},
		}
	}

	client, err := repository.CreateClient(ctx)
	var page Page[customer.State]
	if err != nil {
		return page, err
	}

	defer client.Disconnect(ctx)
	collection := client.
		Database(Database).
		Collection(CustomersCollection)

	totalElements, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return page, err
	}

	page.Items = make([]customer.State, 0)
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
	return page, nil
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
	if err == mongo.ErrNoDocuments {
		return customer.State{}, ErrCustomerNotFound
	}

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
	collection := client.Database(Database).
		Collection(CustomersCollection)

	if state.Id == uuid.Nil {
		state.Id = uuid.New()
		state.CreatedAt = time.Now().UTC()
		state.ModifiedAt = time.Now().UTC()
		_, err = collection.InsertOne(ctx, state)
		return state, err
	}

	lastChange := state.ModifiedAt
	state.ModifiedAt = time.Now().UTC()

	res, err := collection.ReplaceOne(ctx, bson.M{"id": state.Id, "modifiedAt": lastChange}, state)

	if res.ModifiedCount == 0 {
		return state, ErrCustomerConcurrencyIssue
	}

	return state, err
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
