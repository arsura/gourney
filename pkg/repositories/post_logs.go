package repository

import (
	"context"
	"time"

	config "github.com/arsura/gourney/configs"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type PostLogRepositoryProvider interface {
	CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error)
	CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error)
}

type postLogRepository struct {
	postLogCollection *mongo.Collection
	logger            *zap.SugaredLogger
}

func NewPostLogRepository(db *mongo.Client, logger *zap.SugaredLogger, config *config.Config) *postLogRepository {
	postLogCollection := db.Database(config.MongoDB.LogDatabase.Name).Collection(config.MongoDB.LogDatabase.Collections.PostLogs)
	// You could create indexes here!
	return &postLogRepository{postLogCollection, logger}
}

func (r *postLogRepository) CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error) {
	now := time.Now()
	postLog.CreatedAt = now
	postLog.UpdatedAt = now

	result, err := r.postLogCollection.InsertOne(ctx, postLog)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *postLogRepository) CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error) {
	var (
		docs        []interface{}
		insertedIds []primitive.ObjectID

		now = time.Now()
	)
	for _, postLog := range postLogs {
		postLog.CreatedAt = now
		postLog.UpdatedAt = now
		docs = append(docs, postLog)
	}
	result, err := r.postLogCollection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	for _, ids := range result.InsertedIDs {
		insertedIds = append(insertedIds, ids.(primitive.ObjectID))
	}
	return insertedIds, nil
}
