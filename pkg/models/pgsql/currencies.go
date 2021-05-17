package pgsql

import (
	"context"

	"github.com/arsura/moonbase-service/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CurrencyModel struct {
	Pool *pgxpool.Pool
}

func (m *CurrencyModel) Insert(p *models.Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := m.Pool.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (m *CurrencyModel) Get(id int) (*models.Currency, error) {
	return nil, nil
}
