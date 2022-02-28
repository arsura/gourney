package repository_test

import (
	"errors"
	"testing"

	model "github.com/arsura/gourney/pkg/models/pgsql"
	repository "github.com/arsura/gourney/pkg/repositories"
	"github.com/arsura/gourney/pkg/repositories/mocks"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CurrencyRepoTestSuite struct {
	suite.Suite
	mockDbConn *mocks.MockDbConn
	repo       repository.CurrencyRepoProvider
}

func (suite *CurrencyRepoTestSuite) SetupTest() {
	suite.mockDbConn = new(mocks.MockDbConn)
	suite.repo = repository.NewCurrencyRepo(suite.mockDbConn)
}

func (suite *CurrencyRepoTestSuite) Test_Create_Currency_Repo_Success() {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	suite.mockDbConn.On(
		"Exec",
		mock.Anything,
		stmt,
		"RSI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(pgconn.CommandTag("INSERT 0 1"), nil)

	result, err := suite.repo.Create(
		&model.Currency{
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

func (suite *CurrencyRepoTestSuite) Test_Create_Currency_Repo_Failed() {
	stmt := "INSERT INTO currencies(name, amount, total, rise_rate, rise_factor) VALUES($1, $2, $3, $4, $5)"
	suite.mockDbConn.On(
		"Exec",
		mock.Anything,
		stmt,
		"RSI",
		1000.0,
		1000.0,
		0.1,
		10.0,
	).Return(nil, errors.New("failed to insert"))

	result, err := suite.repo.Create(
		&model.Currency{
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

func (suite *CurrencyRepoTestSuite) Test_FindOne_Currency_Repo_Success() {
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	suite.mockDbConn.
		On(
			"QueryRow",
			mock.Anything,
			stmt,
			int64(1),
		).
		Return(suite.mockDbConn).
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

	result, err := suite.repo.FindOneById(1)
	assert.Equal(suite.T(), result, &model.Currency{
		Id:         1,
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyRepoTestSuite) Test_FindOne_Currency_Repo_Failed() {
	stmt := "SELECT id, name, amount, total, rise_rate, rise_factor FROM currencies WHERE id=$1"
	suite.mockDbConn.
		On(
			"QueryRow",
			mock.Anything,
			stmt,
			int64(1),
		).
		Return(suite.mockDbConn).
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

	result, err := suite.repo.FindOneById(1)
	assert.Equal(suite.T(), result, &model.Currency{})
	assert.NotNil(suite.T(), err)
}

func TestCurrencyRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyRepoTestSuite))
}
