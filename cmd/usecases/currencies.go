package usecase

import (
	"context"

	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models/pgsql"
	repository "github.com/arsura/gourney/pkg/repositories"
	"go.uber.org/zap"
)

type CurrencyUsecaseProvider interface {
	Create(ctx context.Context, c *model.Currency) (int64, error)
	FindOneById(ctx context.Context, id int64) (*model.Currency, error)
}

type currencyUsecase struct {
	currencyRepository repository.CurrencyRepoProvider
	logger             *zap.SugaredLogger
}

func NewCurrencyUsecase(repo repository.CurrencyRepoProvider, logger *zap.SugaredLogger) *currencyUsecase {
	return &currencyUsecase{currencyRepository: repo, logger: logger}
}

func (u *currencyUsecase) Create(ctx context.Context, c *model.Currency) (int64, error) {
	var (
		requestId = ctx.Value(constant.RequestIdKey)
	)

	result, err := u.currencyRepository.Create(&model.Currency{
		Name:       c.Name,
		Amount:     c.Amount,
		Total:      c.Total,
		RiseRate:   c.RiseRate,
		RiseFactor: c.RiseFactor,
	})
	if err != nil {
		u.logger.Errorw("failed to create currency", zap.Any("request_id", requestId), zap.Any("error", err))
		return 0, err
	}
	return result, nil
}

func (u *currencyUsecase) FindOneById(ctx context.Context, id int64) (*model.Currency, error) {
	var (
		requestId = ctx.Value(constant.RequestIdKey)
	)

	result, err := u.currencyRepository.FindOneById(int64(id))
	if err != nil {
		u.logger.Errorw("failed to find currency", zap.Any("request_id", requestId), zap.Any("error", err))
		return nil, err
	}
	return result, nil
}
