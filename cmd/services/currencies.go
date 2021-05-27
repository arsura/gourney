package service

import (
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
)

func (s *Service) CreateCurrency(c *pgsql.Currency) (int64, error) {
	result, err := s.Pg.Currencies.Create(&pgsql.Currency{
		Name:       c.Name,
		Amount:     c.Amount,
		Total:      c.Total,
		RiseRate:   c.RiseRate,
		RiseFactor: c.RiseFactor,
	})
	if err != nil {
		s.Logger.Error.Printf("Failed to create currency: %v\n", err)
		return 0, err
	}
	return result, nil
}

func (s *Service) FindCurrency(id int64) (*pgsql.Currency, error) {
	result, err := s.Pg.Currencies.FindOne(int64(id))
	if err != nil {
		s.Logger.Error.Printf("Failed to find currency: %v\n", err)
		return &pgsql.Currency{}, err
	}
	return result, nil
}
