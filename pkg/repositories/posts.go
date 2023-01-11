package repository

import (
	"context"
	"time"

	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	model "github.com/arsura/gourney/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type PostRepositoryProvider interface {
	CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error)
	FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, post *model.Post) (bool, error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error)
}

type postRepository struct {
	postCollection adapter.MongoCollectionProvider
	logger         *zap.SugaredLogger
}

func NewPostRepository(collection *adapter.MongoCollections, logger *zap.SugaredLogger, config *config.Config) *postRepository {
	return &postRepository{collection.PostCollection, logger}
}

func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error) {
	now := time.Now()

	result, err := r.postCollection.InsertOne(ctx, &model.Post{
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *postRepository) FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	result := r.postCollection.FindOne(ctx, bson.M{model.ID: id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	post := &model.Post{}
	err := result.Decode(&post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) UpdatePostById(ctx context.Context, id primitive.ObjectID, post *model.Post) (bool, error) {
	now := time.Now()

	result, err := r.postCollection.UpdateOne(ctx, bson.M{model.ID: id}, bson.M{
		"$set": bson.M{
			model.TITLE:      post.Title,
			model.CONTENT:    post.Content,
			model.UPDATED_AT: now,
		},
	})
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

func (r *postRepository) DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := r.postCollection.DeleteOne(ctx, bson.M{model.ID: id})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
