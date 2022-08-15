package infrastructure

import (
	"happyday/abstract"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProviderSet = wire.NewSet(
	ProviderMongoDbRepository,
	wire.Bind(new(Repository), new(*MongoDbRepository)),
	wire.Bind(new(ReadOnlyRepository), new(*MongoDbRepository)),
)

func ProviderMongoDbRepository(options *options.ClientOptions) *MongoDbRepository {
	return &MongoDbRepository{
		abstract.NewMongoDbRepository(options),
	}
}
