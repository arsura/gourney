package service

import (
	"github.com/arsura/gourney/pkg/models/pgsql"
	"go.uber.org/zap"
)

type CurrencyService struct {
	Logger       *zap.SugaredLogger
	CurrencyRepo pgsql.CurrencyRepoProvider
}

type CurrencyServiceProvider interface {
	Create(c *pgsql.Currency) (int64, error)
	FindOneByID(id int64) (*pgsql.Currency, error)
}

func (s *CurrencyService) Create(c *pgsql.Currency) (int64, error) {
	result, err := s.CurrencyRepo.Create(&pgsql.Currency{
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

func (s *CurrencyService) FindOneByID(id int64) (*pgsql.Currency, error) {
	result, err := s.CurrencyRepo.FindOneByID(int64(id))
	if err != nil {
		s.Logger.Errorf("Failed to find currency: %v", err)
		return nil, err
	}
	return result, nil
}
