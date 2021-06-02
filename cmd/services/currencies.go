package service

import (
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
)

type CurrencyService interface {
	Create(c *pgsql.Currency) (int64, error)
	FindOneById(id int64) (*pgsql.Currency, error)
}

func (s *Service) Create(c *pgsql.Currency) (int64, error) {
	result, err := s.PgRepo.Currencies.Create(&pgsql.Currency{
		Name:       c.Name,
		Amount:     c.Amount,
		Total:      c.Total,
		RiseRate:   c.RiseRate,
		RiseFactor: c.RiseFactor,
	})
	if err != nil {
		s.Logger.Error.Printf("failed to create currency: %v\n", err)
		return 0, err
	}
	return result, nil
}

func (s *Service) FindOneById(id int64) (*pgsql.Currency, error) {
	result, err := s.PgRepo.Currencies.FindOneById(int64(id))
	if err != nil {
		s.Logger.Error.Printf("failed to find currency: %v\n", err)
		return &pgsql.Currency{}, err
	}
	return result, nil
}
