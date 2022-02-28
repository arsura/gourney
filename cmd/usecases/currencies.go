package usecase

import (
	model "github.com/arsura/gourney/pkg/models/pgsql"
	repository "github.com/arsura/gourney/pkg/repositories"
	"go.uber.org/zap"
)

type CurrencyUsecase struct {
	Logger       *zap.SugaredLogger
	CurrencyRepo repository.CurrencyRepoProvider
}

type CurrencyUsecaseProvider interface {
	Create(c *model.Currency) (int64, error)
	FindOneById(id int64) (*model.Currency, error)
}

func (s *CurrencyUsecase) Create(c *model.Currency) (int64, error) {
	result, err := s.CurrencyRepo.Create(&model.Currency{
		Name:       c.Name,
		Amount:     c.Amount,
		Total:      c.Total,
		RiseRate:   c.RiseRate,
		RiseFactor: c.RiseFactor,
	})
	if err != nil {
		s.Logger.Errorf("Failed to create currency: %v", err)
		return 0, err
	}
	return result, nil
}

func (s *CurrencyUsecase) FindOneById(id int64) (*model.Currency, error) {
	result, err := s.CurrencyRepo.FindOneById(int64(id))
	if err != nil {
		s.Logger.Errorf("Failed to find currency: %v", err)
		return nil, err
	}
	return result, nil
}
