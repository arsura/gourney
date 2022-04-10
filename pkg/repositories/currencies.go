package repository

import (
	"context"

	model "github.com/arsura/gourney/pkg/models/pgsql"
)

type CurrencyRepoProvider interface {
	Create(p *model.Currency) (int64, error)
	FindOneById(id int64) (*model.Currency, error)
}

type currencyRepository struct {
	conn DbConn
}

func NewCurrencyRepo(conn DbConn) *currencyRepository {
	return &currencyRepository{conn: conn}
}

func (db *currencyRepository) Create(p *model.Currency) (int64, error) {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	result, err := db.conn.Exec(context.Background(), stmt, p.Name, p.Amount, p.Total, p.RiseRate, p.RiseFactor)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (db *currencyRepository) FindOneById(id int64) (*model.Currency, error) {
	var currency model.Currency
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	err := db.conn.QueryRow(context.Background(), stmt, id).Scan(
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
