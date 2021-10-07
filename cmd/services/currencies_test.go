package service_test

import (
	"errors"
	"testing"

	service "github.com/arsura/gourney/cmd/services"
	"github.com/arsura/gourney/pkg/models/pgsql"
	pgsql_mock "github.com/arsura/gourney/pkg/models/pgsql/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CurrencyServiceTestSuite struct {
	suite.Suite
	mockRepo *pgsql_mock.MockCurrencyRepo
	service  *service.CurrencyService
}

func (suite *CurrencyServiceTestSuite) SetupTest() {
	logger := zaptest.NewLogger(suite.T()).Sugar()
	suite.mockRepo = new(pgsql_mock.MockCurrencyRepo)
	suite.service = &service.CurrencyService{
		Logger:       logger,
		CurrencyRepo: suite.mockRepo,
	}
}

func (suite *CurrencyServiceTestSuite) Test_Create_Currency_Service_Success() {
	suite.mockRepo.On("Create", &pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(1), nil)
	result, err := suite.service.Create(&pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Equal(suite.T(), result, int64(1))
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyServiceTestSuite) Test_Create_Currency_Service_Failed() {
	suite.mockRepo.On("Create", &pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(0), errors.New("Failed to insert"))
	result, err := suite.service.Create(&pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Equal(suite.T(), result, int64(0))
	assert.NotNil(suite.T(), err)
}

func (suite *CurrencyServiceTestSuite) Test_FindOneByID_Currency_Service_Success() {
	suite.mockRepo.On("FindOneByID", int64(1)).Return(&pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}, nil)
	result, err := suite.service.FindOneByID(int64(1))
	assert.Equal(suite.T(), result, &pgsql.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyServiceTestSuite) Test_FindOneByID_Currency_Service_Failed() {
	suite.mockRepo.On("FindOneByID", int64(1)).Return(nil, errors.New("failed to find currency"))
	result, err := suite.service.FindOneByID(int64(1))
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
}

func TestCurrencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}
