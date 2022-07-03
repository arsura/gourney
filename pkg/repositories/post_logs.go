package repository

import (
	"context"
	"time"

	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type PostLogRepositoryProvider interface {
	CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error)
	CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error)
}

type postLogRepository struct {
	postLogCollection adapter.MongoCollectionProvider
	logger            *zap.SugaredLogger
}

func NewPostLogRepository(collection *adapter.MongoCollections, logger *zap.SugaredLogger, config *config.Config) *postLogRepository {
	return &postLogRepository{collection.PostLogCollection, logger}
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
