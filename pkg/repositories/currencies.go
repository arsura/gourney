package repository

import (
	"context"

	model "github.com/arsura/gourney/pkg/models/pgsql"
)

type CurrencyRepoProvider interface {
	Create(p *model.Currency) (int64, error)
	FindOneById(id int64) (*model.Currency, error)
}

type currencyRepo struct {
	Conn DbConn
}

func NewCurrencyRepo(conn DbConn) *currencyRepo {
	return &currencyRepo{Conn: conn}
}

func (db *currencyRepo) Create(p *model.Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := db.Conn.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (db *currencyRepo) FindOneById(id int64) (*model.Currency, error) {
	var currency model.Currency
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
		return &model.Currency{}, err
	}
	return &currency, nil
}
