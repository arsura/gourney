package pgsql

import (
	"testing"

	pgsql_mock "github.com/arsura/moonbase-service/pkg/models/pgsql/mocks"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDB_Insert(t *testing.T) {
	mockDbConn := new(pgsql_mock.MockDBConn)
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	mockDbConn.On(
		"Exec",
		mock.Anything,
		stmt,
		"RSI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(pgconn.CommandTag("INSERT 0 1"), nil)

	db := &DB{Conn: mockDbConn}

	result, err := db.Insert(
		&Currency{
			Name:       "RSI",
			Amount:     1000.0,
			Total:      1000.0,
			RiseRate:   0.1,
			RiseFactor: 10.0,
		},
	)
	assert.Equal(t, result, int64(1))
	assert.Nil(t, err)
}
