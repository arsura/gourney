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

type PostLogUsecaseProvider interface {
	CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error)
	CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error)
}

type postLogUsecase struct {
	repositories *repository.Repository
	services     *service.Services
	logger       *zap.SugaredLogger
}

func NewPostLogUsecase(repositories *repository.Repository, services *service.Services, logger *zap.SugaredLogger) *postLogUsecase {
	return &postLogUsecase{repositories, services, logger}
}

func (u *postLogUsecase) CreatePostLog(ctx context.Context, postLog *model.PostLog) (*primitive.ObjectID, error) {
	id, err := u.repositories.PostLog.CreatePostLog(ctx, &model.PostLog{
		PostId:   postLog.PostId,
		Event:    postLog.Event,
		ActionAt: postLog.ActionAt,
	})
	if err != nil {
		u.logger.With("event", "create_post_log", "error", err, "tracking_id", ctx.Value(constant.REQUEST_ID_KEY)).Error("failed to create post log")
		return nil, err
	}
	return id, nil
}

func (u *postLogUsecase) CreatePostLogs(ctx context.Context, postLogs []model.PostLog) ([]primitive.ObjectID, error) {
	ids, err := u.repositories.PostLog.CreatePostLogs(ctx, postLogs)
	if err != nil {
		u.logger.With("event", "create_post_logs", "error", err).Error("failed to create post logs")
		return nil, err
	}
	return ids, nil
}
