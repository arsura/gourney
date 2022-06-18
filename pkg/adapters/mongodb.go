package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/arsura/gourney/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDBClient(config *config.Config) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		panic(fmt.Errorf("failed to new mongodb client, %v", err))
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to mongodb, %v", err))
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Errorf("failed to ping to mongodb, %v", err))
	}

	return client
}
