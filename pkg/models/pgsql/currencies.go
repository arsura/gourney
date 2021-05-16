package pgsql

import (
	"github.com/arsura/moonbase-service/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CurrencyModel struct {
	Pool *pgxpool.Pool
}

type insertParams struct {
	name       string
	amount     float64
	total      float64
	riseRate   float64
	riseFactor float64
}

func (m *CurrencyModel) Insert(p *insertParams) (string, error) {
	return p.name, nil
}

func (m *CurrencyModel) Get(id int) (*models.Currency, error) {
	return nil, nil
}
