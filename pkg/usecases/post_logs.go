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

type PostLogUseCaseProvider interface {
	CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error)
	CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error)
}

type postLogUseCase struct {
	repositories *repository.Repository
	services     *service.Services
	logger       *zap.SugaredLogger
}

func NewPostLogUseCase(repositories *repository.Repository, services *service.Services, logger *zap.SugaredLogger) *postLogUseCase {
	return &postLogUseCase{repositories, services, logger}
}

func (u *postLogUseCase) CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error) {
	id, err := u.repositories.PostLog.CreatePostLog(ctx, &model.PostLog{
		PostId:   postLog.PostId,
		Event:    postLog.Event,
		ActionAt: postLog.ActionAt,
	})
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to create post log")
		return nil, err
	}
	return id, nil
}

func (u *postLogUseCase) CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error) {
	ids, err := u.repositories.PostLog.CreatePostLogs(ctx, postLogs)
	if err != nil {
		u.logger.With(logger.ERROR_FIELD, err, logger.TRACKING_ID_FIELD, ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to create post logs")
		return nil, err
	}
	return ids, nil
}
