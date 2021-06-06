package handler_test

import (
	"testing"

	"net/http/httptest"

	handler "github.com/arsura/moonbase-service/cmd/handlers"
	service_mock "github.com/arsura/moonbase-service/cmd/services/mocks"
	"github.com/arsura/moonbase-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CurrencyHandlerTestSuite struct {
	suite.Suite
	mockService *service_mock.MockCurrencyServiceProvider
	handler     *handler.CurrencyHandler
	server      *fiber.App
}

func (suite *CurrencyHandlerTestSuite) SetupTest() {
	validate, trans := validator.InitValidate()
	suite.mockService = new(service_mock.MockCurrencyServiceProvider)
	suite.handler = &handler.CurrencyHandler{
		Validator: &validator.Validator{
			Validate: validate,
			Trans:    trans,
		},
		CurrencyService: suite.mockService,
	}
	server := fiber.New()
	server.Post("/", func(c *fiber.Ctx) error {
		return nil
	})
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_Success() {
	server := fiber.New()
	server.Post("/", func(c *fiber.Ctx) error {
		return nil
	})
	_, err := server.Test(httptest.NewRequest(fiber.MethodPost, "/", nil))
	assert.Nil(suite.T(), err)
}

func (suite *CurrencyHandlerTestSuite) Test_Create_Currency_Handler_Failed() {

}

func TestCurrencyHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyHandlerTestSuite))
}
