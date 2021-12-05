package api_handler_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	handler "github.com/arsura/gourney/cmd/api/handlers"
	usecase_mock "github.com/arsura/gourney/cmd/usecases/mocks"
	repo "github.com/arsura/gourney/pkg/repositories"
	"github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CurrencyHandlerTestSuite struct {
	suite.Suite
	mockUsecase *usecase_mock.MockCurrencyUsecaseProvider
	handler     *handler.CurrencyHandler
	server      *fiber.App
}

func (suite *CurrencyHandlerTestSuite) SetupTest() {
	validatr, trans := validator.InitValidator()
	suite.mockUsecase = new(usecase_mock.MockCurrencyUsecaseProvider)
	suite.handler = &handler.CurrencyHandler{
		Validator: &validator.Validator{
			Validate: validatr,
			Trans:    trans,
		},
		CurrencyUsecase: suite.mockUsecase,
	}
	suite.server = fiber.New()
	suite.server.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})
	suite.server.Post("/currencies", suite.handler.CreateCurrencyHandler)
	suite.server.Get("/currencies/:id", suite.handler.FindCurrencyByIdHandler)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_Created() {
	currency := []byte(`{
		"name": "RSI",
		"amount": 1000,
		"total": 1000,
		"riseRate": 0.1,
		"riseFactor": 10.0
	}`)
	request := httptest.NewRequest(fiber.MethodPost, "/currencies", bytes.NewReader(currency))
	request.Header.Set("Content-Type", "application/json")
	suite.mockUsecase.On("Create", &repo.Currency{
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}).Return(int64(1), nil)
	resp, _ := suite.server.Test(request)
	assert.Equal(suite.T(), fiber.StatusCreated, resp.StatusCode)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_BadRequest_Without_Payload() {
	resp, _ := suite.server.Test(httptest.NewRequest(fiber.MethodPost, "/currencies", nil))
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_BadRequest_Invalid_Name() {
	currency := []byte(`{
		"name": 100,
		"amount": 1000,
		"total": 1000,
		"riseRate": 0.1,
		"riseFactor": 10.0
	}`)
	request := httptest.NewRequest(fiber.MethodPost, "/currencies", bytes.NewReader(currency))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(suite.T(), []byte(`{"error":"name must be a string"}`), respBody)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_BadRequest_Missing_Name() {
	currency := []byte(`{
		"amount": 1000,
		"total": 1000,
		"riseRate": 0.1,
		"riseFactor": 10.0
	}`)
	request := httptest.NewRequest(fiber.MethodPost, "/currencies", bytes.NewReader(currency))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(suite.T(), []byte(`{"errors":["name is a required field"]}`), respBody)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_BadRequest_Missing_Name_Amount() {
	currency := []byte(`{
		"total": 1000,
		"riseRate": 0.1,
		"riseFactor": 10.0
	}`)
	request := httptest.NewRequest(fiber.MethodPost, "/currencies", bytes.NewReader(currency))
	request.Header.Set("Content-Type", "application/json")
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(suite.T(), []byte(`{"errors":["name is a required field","amount is a required field"]}`), respBody)
}

func (suite *CurrencyHandlerTestSuite) Test_FindCurrencyById_Handler_Success() {
	request := httptest.NewRequest(fiber.MethodGet, "/currencies/1", nil)
	suite.mockUsecase.On("FindOneById", int64(1)).Return(&repo.Currency{
		Id:         1,
		Name:       "RSI",
		Amount:     1000.0,
		Total:      1000.0,
		RiseRate:   0.1,
		RiseFactor: 10.0,
	}, nil)
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), []byte(`{"id":1,"name":"RSI","amount":1000,"total":1000,"riseRate":0.1,"riseFactor":10}`), respBody)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
}

func (suite *CurrencyHandlerTestSuite) Test_FindCurrencyById_Handler_BadRequest_Invalid_Id() {
	request := httptest.NewRequest(fiber.MethodGet, "/currencies/a", nil)
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), []byte(`{"error":"id must be a number"}`), respBody)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *CurrencyHandlerTestSuite) Test_FindCurrencyById_Handler_NotFound() {
	request := httptest.NewRequest(fiber.MethodGet, "/currencies/1", nil)
	suite.mockUsecase.On("FindOneById", int64(1)).Return(nil, errors.New("failed to find currency"))
	resp, _ := suite.server.Test(request)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), []byte(`{"error":"currency not found"}`), respBody)
	assert.Equal(suite.T(), fiber.StatusNotFound, resp.StatusCode)
}

func TestCurrencyHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyHandlerTestSuite))
}
