package usecase

import (
	"context"

	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	repository "github.com/arsura/gourney/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type PostUsecaseProvider interface {
	CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error)
	FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, update *model.Post) (bool, error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error)
}

type postUsecase struct {
	repositories *repository.Repositories
	logger       *zap.SugaredLogger
}

func NewPostUsecase(repositories *repository.Repositories, logger *zap.SugaredLogger) *postUsecase {
	return &postUsecase{repositories, logger}
}

func (u *postUsecase) CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error) {
	id, err := u.repositories.Posts.CreatePost(ctx, &model.Post{
		Title:   post.Title,
		Content: post.Content,
	})
	if err != nil {
		u.logger.With("error", err, "tracking_id", ctx.Value(constant.RequestIdKey)).Errorf("failed to create post")
		return nil, err
	}

	return id, nil
}

func (u *postUsecase) FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	post, err := u.repositories.Posts.FindPostById(ctx, id)
	if err != nil {
		u.logger.With("error", err, "tracking_id", ctx.Value(constant.RequestIdKey)).Errorf("failed to find post by id")
		return nil, err
	}
	return post, nil
}

func (u *postUsecase) UpdatePostById(ctx context.Context, id primitive.ObjectID, update *model.Post) (bool, error) {
	result, err := u.repositories.Posts.UpdatePostById(ctx, id, &model.Post{
		Title:   update.Title,
		Content: update.Content,
	})
	if err != nil {
		u.logger.With("error", err, "tracking_id", ctx.Value(constant.RequestIdKey)).Errorf("failed to update post by id")
		return false, err
	}
	return result, nil
}

func (u *postUsecase) DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := u.repositories.Posts.DeletePostById(ctx, id)
	if err != nil {
		u.logger.With("error", err, "tracking_id", ctx.Value(constant.RequestIdKey)).Errorf("failed to delete post by id")
		return false, err
	}
	return result, nil
}
