package pgsql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestDB_Insert(t *testing.T) {
	mockDbConn := new(MockDBConn)
	mockDbConn.On(
		"Exec",
		mock.Anything,
		"INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)",
		"SI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(nil, nil)

	db := &DB{Conn: mockDbConn}

	result, err := db.Insert(
		&Currency{
			Name:       "SI",
			Amount:     1000.0,
			Total:      1000.0,
			RiseRate:   0.1,
			RiseFactor: 10.0,
		},
	)
	fmt.Println(result.RowsAffected, err)
	// dbConn.AssertExpectations(t)
}
