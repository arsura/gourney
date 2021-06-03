package pgsql

import (
	"context"
)

type Currency struct {
	ID         int64
	Name       string
	Amount     float64
	Total      float64
	RiseRate   float64
	RiseFactor float64
}

type CurrencyRepo struct {
	Conn DBConn
}

type CurrencyRepoProvider interface {
	Create(p *Currency) (int64, error)
	FindOneByID(id int64) (*Currency, error)
}

func (db *CurrencyRepo) Create(p *Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := db.Conn.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (db *CurrencyRepo) FindOneByID(id int64) (*Currency, error) {
	var currency Currency
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	err := db.Conn.QueryRow(context.Background(), stmt, id).Scan(
		&currency.ID,
		&currency.Name,
		&currency.Amount,
		&currency.Total,
		&currency.RiseRate,
		&currency.RiseFactor,
	)
	if err != nil {
		return &Currency{}, err
	}
	return &currency, nil
}
