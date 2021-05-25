package pgsql

import (
	"errors"
	"testing"

	pgsql_mock "github.com/arsura/moonbase-service/pkg/models/pgsql/mocks"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateCurrencyTestSuite struct {
	suite.Suite
	mockDBConn *pgsql_mock.MockDBConn
}

func (suite *CreateCurrencyTestSuite) SetupTest() {
	suite.mockDBConn = new(pgsql_mock.MockDBConn)
}

func (suite *CreateCurrencyTestSuite) Test_Create_Success() {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	suite.mockDBConn.On(
		"Exec",
		mock.Anything,
		stmt,
		"RSI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(pgconn.CommandTag("INSERT 0 1"), nil)

	db := &DB{Conn: suite.mockDBConn}

	result, err := db.Create(
		&Currency{
			Name:       "RSI",
			Amount:     1000.0,
			Total:      1000.0,
			RiseRate:   0.1,
			RiseFactor: 10.0,
		},
	)
	assert.Equal(suite.T(), result, int64(1))
	assert.Nil(suite.T(), err)
}

func (suite *CreateCurrencyTestSuite) Test_Insert_Failed() {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	suite.mockDBConn.On(
		"Exec",
		mock.Anything,
		stmt,
		"RSI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(nil, errors.New("Failed to insert."))

	db := &DB{Conn: suite.mockDBConn}

	result, err := db.Create(
		&Currency{
			Name:       "RSI",
			Amount:     1000.0,
			Total:      1000.0,
			RiseRate:   0.1,
			RiseFactor: 10.0,
		},
	)
	assert.Equal(suite.T(), result, int64(0))
	assert.NotNil(suite.T(), err)
}

func TestInsertTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCurrencyTestSuite))
}
