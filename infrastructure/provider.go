package infrastructure

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProvideSet = wire.NewSet(
	ProvideProductRepository,

	wire.Bind(new(ProductRepository), new(*MongoDbProductRepository)),
)

func ProvideProductRepository(opt *options.ClientOptions) *MongoDbProductRepository {
	return &MongoDbProductRepository{
		MongoDbRepository{
			options: opt,
		},
	}
}
