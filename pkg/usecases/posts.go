package usecase

import (
	"context"

	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models"
	repository "github.com/arsura/gourney/pkg/repositories"
	service "github.com/arsura/gourney/pkg/services"
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
	repositories *repository.Repository
	services     *service.Services
	logger       *zap.SugaredLogger
}

func NewPostUsecase(repositories *repository.Repository, services *service.Services, logger *zap.SugaredLogger) *postUsecase {
	return &postUsecase{repositories, services, logger}
}

func (u *postUsecase) CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error) {
	u.logger.With("tracking_id", ctx.Value(constant.REQUEST_ID_KEY), "post", post).Info("create new post")

	id, err := u.repositories.Post.CreatePost(ctx, post)
	if err != nil {
		u.logger.With("event", "create_post", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to create post")
		return nil, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_CREATE, &model.Post{Id: *id})
		if err != nil {
			u.logger.With("event", "create_post", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish create post log event")
		}
	}()

	return id, nil
}

func (u *postUsecase) FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	post, err := u.repositories.Post.FindPostById(ctx, id)
	if err != nil {
		u.logger.With("event", "find_post_by_id", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to find post by id")
		return nil, err
	}
	return post, nil
}

func (u *postUsecase) UpdatePostById(ctx context.Context, id primitive.ObjectID, update *model.Post) (bool, error) {
	result, err := u.repositories.Post.UpdatePostById(ctx, id, &model.Post{
		Title:   update.Title,
		Content: update.Content,
	})
	if err != nil {
		u.logger.With("event", "update_post_by_id", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to update post by id")
		return false, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_UPDATE, &model.Post{Id: id})
		if err != nil {
			u.logger.With("event", "update_post_by_id", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish update post log event")
		}
	}()

	return result, nil
}

func (u *postUsecase) DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := u.repositories.Post.DeletePostById(ctx, id)
	if err != nil {
		u.logger.With("event", "delete_post_by_id", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to delete post by id")
		return false, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_DELETE, &model.Post{Id: id})
		if err != nil {
			u.logger.With("event", "delete_post_by_id", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish delete post log event")
		}
	}()

	return result, nil
}
