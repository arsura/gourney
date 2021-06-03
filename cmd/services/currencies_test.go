package service_test

import (
	"errors"
	"testing"

	service "github.com/arsura/moonbase-service/cmd/services"
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	pgsql_mock "github.com/arsura/moonbase-service/pkg/models/pgsql/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CurrencyTestSuite struct {
	suite.Suite
	mockRepositories *pgsql_mock.MockCurrencyRepository
	mockService      *service.Service
}

func (suite *CurrencyTestSuite) SetupTest() {
	suite.mockRepositories = new(pgsql_mock.MockCurrencyRepository)
	suite.mockService = &service.Service{
		PgRepo: &pgsql.Repositories{
			Currencies: suite.mockRepositories,
		},
	}
}

func (suite *CurrencyTestSuite) Test_Create_Success() {
	suite.mockRepositories.On("Create", &pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(1), nil)

	result, err := suite.mockService.Create(&pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Equal(suite.T(), result, int64(1))
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyTestSuite) Test_Create_Failed() {
	suite.mockRepositories.On("Create", &pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(0), errors.New("failed to insert"))

	result, err := suite.mockService.Create(&pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})

	assert.Equal(suite.T(), result, int64(0))
	assert.NotNil(suite.T(), err)
}

func TestCurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyTestSuite))
}
