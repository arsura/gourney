package pgsql

import (
	"context"
)

type Currency struct {
	ID         int
	Name       string
	Amount     float64
	Total      float64
	RiseRate   float64
	RiseFactor float64
}

type CurrencyRepository interface {
	Insert(p *Currency) (int64, error)
	Get(id int) (*Currency, error)
}

func (db *DB) Insert(p *Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := db.Conn.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (m *DB) Get(id int) (*Currency, error) {
	return nil, nil
}
