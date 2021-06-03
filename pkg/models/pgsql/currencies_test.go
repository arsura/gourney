package pgsql_test

import (
	"errors"
	"testing"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	pgsql_mock "github.com/arsura/moonbase-service/pkg/models/pgsql/mocks"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CurrencyTestSuite struct {
	suite.Suite
	mockDBConn *pgsql_mock.MockDBConn
	mockDB     *pgsql.DB
}

func (suite *CurrencyTestSuite) SetupTest() {
	suite.mockDBConn = new(pgsql_mock.MockDBConn)
	suite.mockDB = &pgsql.DB{Conn: suite.mockDBConn}
}

func (suite *CurrencyTestSuite) Test_Create_Success() {
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

	result, err := suite.mockDB.Create(
		&pgsql.Currency{
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

func (suite *CurrencyTestSuite) Test_Insert_Failed() {
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
	).Return(nil, errors.New("failed to insert"))

	result, err := suite.mockDB.Create(
		&pgsql.Currency{
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

func (suite *CurrencyTestSuite) Test_FindOne_Success() {
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	suite.mockDBConn.
		On(
			"QueryRow",
			mock.Anything,
			stmt,
			int64(1),
		).
		Return(suite.mockDBConn).
		On(
			"Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Run(func(args mock.Arguments) {
			arg1 := args.Get(0).(*int64)
			*arg1 = 1
			arg2 := args.Get(1).(*string)
			*arg2 = "RSI"
			arg3 := args.Get(2).(*float64)
			*arg3 = 1000.0
			arg4 := args.Get(3).(*float64)
			*arg4 = 1000.0
			arg5 := args.Get(4).(*float64)
			*arg5 = 0.1
			arg6 := args.Get(5).(*float64)
			*arg6 = 10.0
		}).
		Return(nil)

	result, err := suite.mockDB.FindOneById(1)
	assert.Equal(suite.T(), result, &pgsql.Currency{
		ID:         1,
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyTestSuite) Test_FindOne_Failed() {
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	suite.mockDBConn.
		On(
			"QueryRow",
			mock.Anything,
			stmt,
			int64(1),
		).
		Return(suite.mockDBConn).
		On(
			"Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(errors.New("failed to query"))

	result, err := suite.mockDB.FindOneById(1)
	assert.Equal(suite.T(), result, &pgsql.Currency{})
	assert.NotNil(suite.T(), err)
}

func TestCurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyTestSuite))
}
