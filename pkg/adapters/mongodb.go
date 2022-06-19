package adapter

import (
	"context"
	"time"

	config "github.com/arsura/gourney/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func NewMongoDBClient(logger *zap.SugaredLogger, config *config.Config) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		logger.With("error", err).Panic("failed to new mongodb client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.With("error", err).Panic("failed to connect to mongodb")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.With("error", err).Panic("failed to ping to mongodb")
	}

	return client
}
