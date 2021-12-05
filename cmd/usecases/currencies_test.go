package usecase_test

import (
	"errors"
	"testing"

	usecase "github.com/arsura/gourney/cmd/usecases"
	repo "github.com/arsura/gourney/pkg/repositories"
	pgsql_mock "github.com/arsura/gourney/pkg/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CurrencyUsecaseTestSuite struct {
	suite.Suite
	mockRepo *pgsql_mock.MockCurrencyRepo
	usecase  *usecase.CurrencyUsecase
}

func (suite *CurrencyUsecaseTestSuite) SetupTest() {
	logger := zaptest.NewLogger(suite.T()).Sugar()
	suite.mockRepo = new(pgsql_mock.MockCurrencyRepo)
	suite.usecase = &usecase.CurrencyUsecase{
		Logger:       logger,
		CurrencyRepo: suite.mockRepo,
	}
}

func (suite *CurrencyUsecaseTestSuite) Test_Create_Currency_Usecase_Success() {
	suite.mockRepo.On("Create", &repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(1), nil)
	result, err := suite.usecase.Create(&repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Equal(suite.T(), result, int64(1))
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyUsecaseTestSuite) Test_Create_Currency_Usecase_Failed() {
	suite.mockRepo.On("Create", &repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(0), errors.New("Failed to insert"))
	result, err := suite.usecase.Create(&repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Equal(suite.T(), result, int64(0))
	assert.NotNil(suite.T(), err)
}

func (suite *CurrencyUsecaseTestSuite) Test_FindOneById_Currency_Usecase_Success() {
	suite.mockRepo.On("FindOneById", int64(1)).Return(&repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}, nil)
	result, err := suite.usecase.FindOneById(int64(1))
	assert.Equal(suite.T(), result, &repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	})
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyUsecaseTestSuite) Test_FindOneById_Currency_Usecase_Failed() {
	suite.mockRepo.On("FindOneById", int64(1)).Return(nil, errors.New("failed to find currency"))
	result, err := suite.usecase.FindOneById(int64(1))
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
}

func TestCurrencyUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyUsecaseTestSuite))
}
