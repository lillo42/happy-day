package infrastructure

import (
	"context"
	"math"
	"time"

	"happyday/abstract"
	"happyday/common"
	"happyday/product/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ID      = "id"
	Version = "version"
	Name    = "name"
	Price   = "price"

	NameAsc   OrderBy = "nameAsc"
	NameDesc  OrderBy = "nameDesc"
	PriceAsc  OrderBy = "priceAsc"
	PriceDesc OrderBy = "priceDesc"

	Collection = "products"
)

var (
	NotFound         = common.NewError("CLR000", "Club not found")
	ConcurrencyIssue = common.NewError("CLR001", "Concurrency issue")

	_ Repository         = (*MongoDbRepository)(nil)
	_ ReadOnlyRepository = (*MongoDbRepository)(nil)
)

type (
	Repository interface {
		Create(id uuid.UUID) domain.AggregateRoot
		Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error)
		Save(ctx context.Context, root domain.AggregateRoot) error
		Delete(ctx context.Context, root domain.AggregateRoot) error
	}

	ReadOnlyRepository interface {
		GetById(context context.Context, id uuid.UUID) (DetailsViewModel, error)
		GetAll(context context.Context, filter Filter) (common.Page[ViewModel], error)
		Exists(ctx context.Context, product uuid.UUID) (bool, error)
	}

	DetailsViewModel struct {
		ID       uuid.UUID      `bson:"id"`
		Name     string         `bson:"name"`
		Price    float64        `bson:"price"`
		IsEnable bool           `bson:"isEnable"`
		Priority int64          `bson:"priority"`
		Products []InnerProduct `bson:"products"`
		CreateAt time.Time      `bson:"createAt"`
		ModifyAt time.Time      `bson:"modifyAt"`
	}

	ViewModel struct {
		ID       uuid.UUID
		Name     string
		Price    float64
		IsEnable bool
	}

	OrderBy string

	Filter struct {
		OrderBy OrderBy
		Text    string
		Page    int64
		Size    int64
	}

	MongoDbRepository struct {
		abstract.MongoDbRepository
	}

	Entity struct {
		ID       uuid.UUID      `bson:"id"`
		Name     string         `bson:"name"`
		Price    float64        `bson:"price"`
		Priority int64          `bson:"priority"`
		Products []InnerProduct `bson:"products"`
		IsEnable bool           `bson:"isEnable"`
		CreateAt time.Time      `bson:"createAt"`
		ModifyAt time.Time      `bson:"modifyAt"`
		Version  int64          `bson:"version"`
	}

	InnerProduct struct {
		ID     uuid.UUID `bson:"id"`
		Amount int64     `bson:"amount"`
	}
)

func (repository MongoDbRepository) Exists(ctx context.Context, product uuid.UUID) (bool, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return false, err
	}

	defer client.Disconnect(ctx)
	totalElements, err := client.Database(common.Database).
		Collection(Collection).
		CountDocuments(ctx, bson.M{ID: product})

	if err != nil {
		return false, err
	}

	return totalElements == 1, nil
}

func (repository MongoDbRepository) GetAll(ctx context.Context, filter Filter) (common.Page[ViewModel], error) {
	opt := options.Find().
		SetLimit(filter.Size).
		SetSkip(filter.Size * filter.Page)

	if filter.OrderBy == NameAsc {
		opt = opt.SetSort(bson.M{Name: 1})
	} else if filter.OrderBy == NameDesc {
		opt = opt.SetSort(bson.M{Name: -1})
	} else if filter.OrderBy == PriceAsc {
		opt = opt.SetSort(bson.M{Price: 1})
	} else if filter.OrderBy == PriceDesc {
		opt = opt.SetSort(bson.M{Price: -1})
	}

	query := bson.M{}
	if len(filter.Text) > 0 {
		query["$text"] = bson.M{"$search": filter.Text}
	}

	client, err := repository.CreateClient(ctx)
	if err != nil {
		return common.Page[ViewModel]{}, err
	}

	defer client.Disconnect(ctx)

	page := common.Page[ViewModel]{}
	collection := client.Database(common.Database).
		Collection(Collection)

	totalElements, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return page, err
	}

	cursor, err := collection.
		Find(ctx, query, opt)

	if err != nil {
		return page, err
	}

	page.Items = make([]ViewModel, 0)
	page.TotalElements = totalElements

	if totalElements > 0 {
		tmp := float64(totalElements) / float64(filter.Size)
		tmp = math.Ceil(tmp)
		page.TotalPages = int64(tmp)
	}

	for cursor.TryNext(ctx) {
		err = cursor.Err()
		if err != nil {
			return page, err
		}

		var viewModel ViewModel
		cursor.Decode(&viewModel)
		page.Items = append(page.Items, viewModel)
	}

	return page, nil
}

func (repository MongoDbRepository) GetById(ctx context.Context, id uuid.UUID) (DetailsViewModel, error) {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return DetailsViewModel{}, err
	}

	defer client.Disconnect(ctx)

	var viewModel DetailsViewModel
	err = client.Database(common.Database).
		Collection(Collection).
		FindOne(ctx, bson.M{ID: id}).
		Decode(&viewModel)

	return viewModel, nil
}

func (repository MongoDbRepository) Create(id uuid.UUID) domain.AggregateRoot {
	return domain.NewAggregateRoot(domain.NewStateWithID(id), 0)
}

func (repository MongoDbRepository) Get(ctx context.Context, id uuid.UUID) (domain.AggregateRoot, error) {
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

	products := make([]domain.Product, len(entity.Products))

	for index, product := range entity.Products {
		products[index] = domain.NewProduct(product.ID, product.Amount)
	}

	return domain.NewAggregateRoot(
			domain.NewState(entity.ID,
				entity.Name,
				entity.Price,
				entity.IsEnable,
				entity.Priority,
				products),
			entity.Version),
		nil
}

func (repository MongoDbRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	if len(root.Events()) == 0 {
		return nil
	}

	state := root.State()

	products := make([]InnerProduct, len(state.Products()))
	for index, product := range state.Products() {
		products[index] = InnerProduct{
			ID:     product.ID(),
			Amount: product.Amount(),
		}
	}

	entity := Entity{
		ID:       state.ID(),
		Name:     state.Name(),
		Price:    state.Price(),
		Priority: state.Priority(),
		Products: products,
		CreateAt: time.Now().UTC(),
		ModifyAt: time.Now().UTC(),
		Version:  root.Version() + int64(len(root.Events())),
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

func (repository MongoDbRepository) Delete(ctx context.Context, root domain.AggregateRoot) error {
	client, err := repository.CreateClient(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)

	_, err = client.Database(common.Database).
		Collection(Collection).
		DeleteOne(ctx, bson.M{ID: root.State().ID()})

	return err
}
