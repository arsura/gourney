package repository

import (
	"context"
	"time"

	config "github.com/arsura/gourney/configs"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type PostRepositoryProvider interface {
	CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error)
	FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, post *model.Post) (bool, error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error)
}

type postRepository struct {
	postCollection *mongo.Collection
	logger         *zap.SugaredLogger
}

func NewPostRepository(db *mongo.Client, logger *zap.SugaredLogger, config *config.Config) *postRepository {
	postCollection := db.Database(config.MongoDB.BlogDatabase.Name).Collection(config.MongoDB.BlogDatabase.Collections.Posts)
	// You could create indexes here!
	return &postRepository{postCollection, logger}
}

func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error) {
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	result, err := r.postCollection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *postRepository) FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	result := r.postCollection.FindOne(ctx, bson.M{"_id": id})
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

	result, err := r.postCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": now,
		},
	})
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

func (r *postRepository) DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := r.postCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
