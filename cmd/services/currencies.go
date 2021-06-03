package service

import (
	logger "github.com/arsura/moonbase-service/pkg/logger"
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
)

type CurrencyService struct {
	Logger       *logger.Logger
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
		// s.Logger.Error.Printf("failed to create currency: %v\n", err)
		return 0, err
	}
	return result, nil
}

func (s *CurrencyService) FindOneByID(id int64) (*pgsql.Currency, error) {
	result, err := s.CurrencyRepo.FindOneByID(int64(id))
	if err != nil {
		s.Logger.Error.Printf("failed to find currency: %v\n", err)
		return &pgsql.Currency{}, err
	}
	return result, nil
}
