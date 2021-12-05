package repository

import (
	"context"
)

type Currency struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Total      float64 `json:"total"`
	RiseRate   float64 `json:"riseRate"`
	RiseFactor float64 `json:"riseFactor"`
}

type CurrencyRepo struct {
	Conn DbConn
}

type CurrencyRepoProvider interface {
	Create(p *Currency) (int64, error)
	FindOneById(id int64) (*Currency, error)
}

func (db *CurrencyRepo) Create(p *Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := db.Conn.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (db *CurrencyRepo) FindOneById(id int64) (*Currency, error) {
	var currency Currency
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	err := db.Conn.QueryRow(context.Background(), stmt, id).Scan(
		&currency.Id,
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
