package usecase

import (
	model "github.com/arsura/gourney/pkg/models/pgsql"
	repository "github.com/arsura/gourney/pkg/repositories"
	"go.uber.org/zap"
)

type CurrencyUsecaseProvider interface {
	Create(c *model.Currency) (int64, error)
	FindOneById(id int64) (*model.Currency, error)
}

type currencyUsecase struct {
	CurrencyRepo repository.CurrencyRepoProvider
	Logger       *zap.SugaredLogger
}

func NewCurrencyUsecase(repo repository.CurrencyRepoProvider, logger *zap.SugaredLogger) *currencyUsecase {
	return &currencyUsecase{CurrencyRepo: repo, Logger: logger}
}

func (s *currencyUsecase) Create(c *model.Currency) (int64, error) {
	result, err := s.CurrencyRepo.Create(&model.Currency{
		Name:       c.Name,
		Amount:     c.Amount,
		Total:      c.Total,
		RiseRate:   c.RiseRate,
		RiseFactor: c.RiseFactor,
	})
	if err != nil {
		s.Logger.Errorf("failed to create currency: %v", err)
		return 0, err
	}
	return result, nil
}

func (s *currencyUsecase) FindOneById(id int64) (*model.Currency, error) {
	result, err := s.CurrencyRepo.FindOneById(int64(id))
	if err != nil {
		s.Logger.Errorf("failed to find currency: %v", err)
		return nil, err
	}
	return result, nil
}
