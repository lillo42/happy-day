package infrastructure

import (
	"context"
	"math"
	"time"

	"happyday/abstract"
	"happyday/common"
	"happyday/customer/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Repository interface {
		Create(id uuid.UUID) domain.AggregateRoot
		Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error)
		Save(ctx context.Context, root domain.AggregateRoot) error
		Delete(ctx context.Context, root domain.AggregateRoot) error
	}

	DetailsViewModel struct {
		ID       uuid.UUID `json:"id,omitempty"`
		Name     string    `json:"name,omitempty"`
		Comment  string    `json:"comment,omitempty"`
		Phones   []string  `json:"phones,omitempty"`
		CreateAt time.Time `json:"createAt"`
		ModifyAt time.Time `json:"modifyAt"`
	}

	ViewModel struct {
		ID      uuid.UUID `json:"id,omitempty"`
		Name    string    `json:"name,omitempty"`
		Comment string    `json:"comment,omitempty"`
		Phones  []string  `json:"phones,omitempty"`
	}

	OrderBy string
	Filter  struct {
		OrderBy OrderBy
		Text    string
		Page    int64
		Size    int64
	}

	ReadOnlyRepository interface {
		GetAll(ctx context.Context, filter Filter) (common.Page[ViewModel], error)
		GetById(ctx context.Context, id uuid.UUID) (DetailsViewModel, error)
	}

	MongoDbRepository struct {
		abstract.MongoDbRepository
	}

	Entity struct {
		ID       uuid.UUID `json:"id,omitempty"`
		Name     string    `json:"name,omitempty"`
		Comment  string    `json:"comment,omitempty"`
		Phones   []string  `json:"phones,omitempty"`
		CreateAt time.Time `json:"createAt"`
		ModifyAt time.Time `json:"modifyAt"`
		Version  int64     `json:"version,omitempty"`
	}
)

const (
	NameAsc     OrderBy = "nameAsc"
	NameDesc    OrderBy = "nameDesc"
	CommentAsc  OrderBy = "commentAsc"
	CommentDesc OrderBy = "commentDesc"

	ID         = "id"
	Name       = "name"
	Comment    = "comment"
	Version    = "version"
	Collection = "customers"
)

var (
	NotFound         = common.NewError("CSR000", "Customer not found")
	ConcurrencyIssue = common.NewError("CSR001", "Concurrency issue")
)

func (repository *MongoDbRepository) GetAll(ctx context.Context, filter Filter) (common.Page[ViewModel], error) {
	opt := options.Find().
		SetLimit(filter.Size).
		SetSkip(filter.Size * filter.Page)

	if filter.OrderBy == NameAsc {
		opt = opt.SetSort(bson.M{Name: 1})
	} else if filter.OrderBy == NameDesc {
		opt = opt.SetSort(bson.M{Name: -1})
	} else if filter.OrderBy == CommentAsc {
		opt = opt.SetSort(bson.M{Comment: 1})
	} else if filter.OrderBy == CommentDesc {
		opt = opt.SetSort(bson.M{Comment: -1})
	}

	query := bson.M{}
	if len(filter.Text) > 0 {
		query["$text"] = bson.M{"$search": filter.Text}
	}

	var res common.Page[ViewModel]
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return res, err
	}

	defer client.Disconnect(ctx)
	collection := client.Database(common.Database).
		Collection(Collection)

	totalElements, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return res, err
	}

	cursor, err := collection.
		Find(ctx, query, opt)

	if err != nil {
		return res, err
	}

	res.Items = make([]ViewModel, 0)
	res.TotalElements = totalElements

	if totalElements > 0 {
		tmp := float64(totalElements) / float64(filter.Size)
		tmp = math.Ceil(tmp)
		res.TotalPages = int64(tmp)
	}

	for cursor.TryNext(ctx) {
		err = cursor.Err()
		if err != nil {
			return res, err
		}

		var viewModel ViewModel
		cursor.Decode(&viewModel)
		res.Items = append(res.Items, viewModel)
	}

	return res, nil
}

func (repository *MongoDbRepository) GetById(ctx context.Context, id uuid.UUID) (DetailsViewModel, error) {
	var res DetailsViewModel
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return res, err
	}

	defer client.Disconnect(ctx)
	err = client.Database(common.Database).
		Collection(Collection).
		FindOne(ctx, bson.M{ID: id}).
		Decode(&res)

	if err == mongo.ErrNoDocuments {
		return res, NotFound
	} else if err != nil {
		return res, err
	}
	return res, nil
}

func (repository *MongoDbRepository) Create(id uuid.UUID) domain.AggregateRoot {
	return domain.NewAggregateRoot(domain.NewStateWithID(id), 0)
}

func (repository *MongoDbRepository) Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)
	var entity Entity
	err = client.Database(common.Database).
		Collection(Collection).
		FindOne(ctx, bson.M{ID: id}).
		Decode(&entity)

	if err == mongo.ErrNoDocuments {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	phones := make([]domain.Phone, len(entity.Phones))
	for index, phone := range entity.Phones {
		phones[index] = domain.NewPhone(phone)
	}

	return domain.NewAggregateRoot(
		domain.NewState(entity.ID,
			entity.Name,
			entity.Comment,
			phones),
		entity.Version), nil
}

func (repository *MongoDbRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	state := root.State()
	entity := Entity{
		ID:      state.ID(),
		Name:    state.Name(),
		Comment: state.Comment(),
		Phones:  make([]string, len(state.Phones())),
	}

	for index, phone := range state.Phones() {
		entity.Phones[index] = phone.Number()
	}

	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)
	collection := client.Database(common.Database).
		Collection(Collection)

	if root.Version() == 0 {
		_, err = collection.InsertOne(ctx, entity)
		return err
	}

	var oldEntity Entity
	err = collection.
		FindOne(ctx, bson.M{ID: state.ID(), Version: root.Version()}).
		Decode(&oldEntity)

	if err == mongo.ErrNoDocuments {
		return NotFound
	} else if err != nil {
		return err
	}

	entity.CreateAt = oldEntity.CreateAt
	res, err := collection.ReplaceOne(ctx, bson.M{ID: state.ID(), Version: root.Version()}, entity)

	if res.ModifiedCount == 0 {
		return ConcurrencyIssue
	}

	return err
}

func (repository *MongoDbRepository) Delete(ctx context.Context, root domain.AggregateRoot) error {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)
	_, err = client.Database(common.Database).
		Collection(Collection).
		DeleteOne(ctx, bson.M{ID: root.State().ID()})

	if err == mongo.ErrNoDocuments {
		return NotFound
	}
	return err
}
