package usecase

import (
	"context"

	"github.com/arsura/gourney/pkg/constant"
	"github.com/arsura/gourney/pkg/logger"
	model "github.com/arsura/gourney/pkg/models"
	repository "github.com/arsura/gourney/pkg/repositories"
	service "github.com/arsura/gourney/pkg/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type PostUseCaseProvider interface {
	CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error)
	FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, update *model.Post) (bool, error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error)
	CountPostBySocialNetworkType(ctx context.Context) ([]repository.CountPostBySocialNetworkTypeRes, error)
}

type postUseCase struct {
	repositories *repository.Repository
	services     *service.Services
	logger       *zap.SugaredLogger
}

func NewPostUseCase(repositories *repository.Repository, services *service.Services, logger *zap.SugaredLogger) *postUseCase {
	return &postUseCase{repositories, services, logger}
}

func (u *postUseCase) CreatePost(ctx context.Context, post *model.Post) (*primitive.ObjectID, error) {
	u.logger.With(logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY), logger.DATA_FIELD, map[string]interface{}{
		"post": post,
	}).Info("create new post")

	id, err := u.repositories.Post.CreatePost(ctx, post)
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to create post")
		return nil, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_CREATE, &model.Post{Id: *id})
		if err != nil {
			u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish create post log event")
		}
	}()

	return id, nil
}

func (u *postUseCase) FindPostById(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	post, err := u.repositories.Post.FindPostById(ctx, id)
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to find post by id")
		return nil, err
	}
	return post, nil
}

func (u *postUseCase) UpdatePostById(ctx context.Context, id primitive.ObjectID, update *model.Post) (bool, error) {
	result, err := u.repositories.Post.UpdatePostById(ctx, id, &model.Post{
		Title:             update.Title,
		Content:           update.Content,
		SocialNetworkType: update.SocialNetworkType,
	})
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to update post by id")
		return false, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_UPDATE, &model.Post{Id: id})
		if err != nil {
			u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish update post log event")
		}
	}()

	return result, nil
}

func (u *postUseCase) DeletePostById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := u.repositories.Post.DeletePostById(ctx, id)
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to delete post by id")
		return false, err
	}

	go func() {
		err = u.services.Log.PublishPostLogEvent(ctx, constant.ACTION_DELETE, &model.Post{Id: id})
		if err != nil {
			u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to publish delete post log event")
		}
	}()

	return result, nil
}

func (u *postUseCase) CountPostBySocialNetworkType(ctx context.Context) ([]repository.CountPostBySocialNetworkTypeRes, error) {
	count, err := u.repositories.Post.CountPostBySocialNetworkType(ctx)
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to find post by id")
		return nil, err
	}
	return count, nil
}
