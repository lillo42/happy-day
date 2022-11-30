package infrastructure

import (
	"context"
	"errors"
	"happy_day/domain/reservation"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ReservationCollection = "reservations"

	DeliveryAsc  ReservationOrderBy = "delivery_asc"
	DeliveryDesc ReservationOrderBy = "delivery_desc"
	PickupAsc    ReservationOrderBy = "pickup_asc"
	PickupDesc   ReservationOrderBy = "pickup_desc"
)

var (
	ErrReservationConcurrencyIssue = errors.New("reservation concurrency issue")
	ErrReservationNotFound         = errors.New("reservation not found")
)

type (
	ReservationOrderBy string
	ReservationFilter  struct {
		Text   string
		Page   int64
		Size   int64
		SortBy ReservationOrderBy
	}

	ReservationRepository interface {
		GetAll(ctx context.Context, filter ReservationFilter) (Page[reservation.State], error)
		GetById(ctx context.Context, id uuid.UUID) (reservation.State, error)
		Save(ctx context.Context, state reservation.State) (reservation.State, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	MockReservationRepository struct {
		mock.Mock
	}

	MongoDbReservationRepository struct {
		MongoDbRepository
	}
)

func (repository MongoDbReservationRepository) GetAll(ctx context.Context, filter ReservationFilter) (Page[reservation.State], error) {
	opt := options.Find().
		SetSkip((filter.Page - 1) * filter.Size).
		SetLimit(filter.Size)

	if filter.SortBy == DeliveryAsc {
		opt.SetSort(bson.D{{"delivery.at", 1}})
	} else if filter.SortBy == DeliveryDesc {
		opt.SetSort(bson.D{{"delivery.at", -1}})
	} else if filter.SortBy == PickupAsc {
		opt.SetSort(bson.D{{"pickUp.at", 1}})
	} else if filter.SortBy == PickupDesc {
		opt.SetSort(bson.D{{"pickUp.at", -1}})
	}

	query := bson.M{}
	if len(filter.Text) > 0 {
		query["$or"] = []interface{}{
			bson.M{"id": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"comment": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.street": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.number": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.neighborhood": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.complement": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.postalCode": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"address.city": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"customer.name": bson.M{"$regex": filter.Text, "$options": "im"}},
			bson.M{"customer.phones": bson.M{"$elemMatch": bson.M{"number": bson.M{"$regex": filter.Text, "$options": "im"}}}},
		}
	}

	client, err := repository.CreateClient(ctx)
	var page Page[reservation.State]
	if err != nil {
		return page, err
	}

	collection := client.
		Database(Database).
		Collection(ReservationCollection)

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

func (repository MongoDbReservationRepository) GetById(ctx context.Context, id uuid.UUID) (reservation.State, error) {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return reservation.State{}, err
	}

	defer client.Disconnect(ctx)
	decode := client.Database(Database).
		Collection(ReservationCollection).
		FindOne(ctx, query)

	err = decode.Err()
	if err == mongo.ErrNoDocuments {
		return reservation.State{}, ErrProductNotFound
	}
	if err != nil {
		return reservation.State{}, err
	}

	var state reservation.State
	err = decode.Decode(&state)
	return state, nil
}

func (repository MongoDbReservationRepository) Save(ctx context.Context, state reservation.State) (reservation.State, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return reservation.State{}, err
	}

	defer client.Disconnect(ctx)
	collection := client.Database(Database).
		Collection(ReservationCollection)

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
		return state, ErrReservationConcurrencyIssue
	}

	return state, err
}

func (repository MongoDbReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := bson.M{"id": id}
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(Database).
		Collection(ReservationCollection).
		DeleteOne(ctx, query)
	return err
}

func (m *MockReservationRepository) GetAll(ctx context.Context, filter ReservationFilter) (Page[reservation.State], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(Page[reservation.State]), args.Error(1)
}

func (m *MockReservationRepository) GetById(ctx context.Context, id uuid.UUID) (reservation.State, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(reservation.State), args.Error(1)
}

func (m *MockReservationRepository) Save(ctx context.Context, state reservation.State) (reservation.State, error) {
	args := m.Called(ctx, state)
	return args.Get(0).(reservation.State), args.Error(1)
}

func (m *MockReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
