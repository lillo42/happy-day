package infrastructure

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProvideSet = wire.NewSet(
	ProvideProductRepository,
	wire.Bind(new(ProductRepository), new(*MongoDbProductRepository)),

	ProviderCustomerRepository,
	wire.Bind(new(CustomerRepository), new(*MongoDbCustomerRepository)),

	ProvideReservationRepository,
	wire.Bind(new(ReservationRepository), new(*MongoDbReservationRepository)),
)

func ProvideProductRepository(opt *options.ClientOptions) *MongoDbProductRepository {
	return &MongoDbProductRepository{
		MongoDbRepository{
			options: opt,
		},
	}
}

func ProviderCustomerRepository(opt *options.ClientOptions) *MongoDbCustomerRepository {
	return &MongoDbCustomerRepository{
		MongoDbRepository{
			options: opt,
		},
	}
}

func ProvideReservationRepository(opt *options.ClientOptions) *MongoDbReservationRepository {
	return &MongoDbReservationRepository{
		MongoDbRepository{
			options: opt,
		},
	}
}
