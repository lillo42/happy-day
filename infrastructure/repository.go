package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbRepository struct {
	options *options.ClientOptions
}

func (repository MongoDbRepository) CreateClient(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, repository.options)
}

func NewMongoDbRepository(options *options.ClientOptions) MongoDbRepository {
	return MongoDbRepository{
		options: options,
	}
}
