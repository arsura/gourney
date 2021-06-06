package handler

import (
	"strconv"

	service "github.com/arsura/moonbase-service/cmd/services"
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	util "github.com/arsura/moonbase-service/pkg/utils"
	validator "github.com/arsura/moonbase-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type CreateReqBody struct {
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Total      float64 `json:"total" validate:"required"`
	RiseRate   float64 `json:"riseRate"`
	RiseFactor float64 `json:"riseFactor"`
}

type CurrencyHandler struct {
	Validator       *validator.Validator
	CurrencyService service.CurrencyServiceProvider
}

type CurrencyHandlerProvider interface {
	CreateCurrencyHandler(c *fiber.Ctx) error
	FindCurrencyByIDHandler(c *fiber.Ctx) error
}

func (h *CurrencyHandler) CreateCurrencyHandler(c *fiber.Ctx) error {
	currency := new(CreateReqBody)
	if err := c.BodyParser(currency); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": util.UnmarshalErrorParser(err),
		})
	}

	if err := h.Validator.Validate.Struct(*currency); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"errors": h.Validator.TransError(err),
		})
	}

	_, err := h.CurrencyService.Create(&pgsql.Currency{
		Name:       currency.Name,
		Amount:     currency.Amount,
		Total:      currency.Total,
		RiseRate:   currency.RiseRate,
		RiseFactor: currency.RiseFactor,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to create currency.",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *CurrencyHandler) FindCurrencyByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be a number.",
		})
	}

	result, err := h.CurrencyService.FindOneByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "currency not found.",
		})

	}
	return c.JSON(result)
}
