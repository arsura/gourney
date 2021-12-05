package usecase

import (
	"github.com/arsura/gourney/pkg/repositories"
	"go.uber.org/zap"
)

type CurrencyUsecase struct {
	Logger       *zap.SugaredLogger
	CurrencyRepo repository.CurrencyRepoProvider
}

type CurrencyUsecaseProvider interface {
	Create(c *repository.Currency) (int64, error)
	FindOneById(id int64) (*repository.Currency, error)
}

func (s *CurrencyUsecase) Create(c *repository.Currency) (int64, error) {
	result, err := s.CurrencyRepo.Create(&repository.Currency{
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

func (s *CurrencyUsecase) FindOneById(id int64) (*repository.Currency, error) {
	result, err := s.CurrencyRepo.FindOneById(int64(id))
	if err != nil {
		s.Logger.Errorf("Failed to find currency: %v", err)
		return nil, err
	}
	return result, nil
}
