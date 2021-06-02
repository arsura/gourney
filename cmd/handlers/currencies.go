package handler

import (
	"strconv"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	"github.com/arsura/moonbase-service/pkg/util"
	"github.com/gofiber/fiber/v2"
)

type CreateReqBody struct {
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Total      float64 `json:"total" validate:"required"`
	RiseRate   float64 `json:"riseRate"`
	RiseFactor float64 `json:"riseFactor"`
}

func (h *Handler) CreateCurrencyHandler(c *fiber.Ctx) error {
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

	_, err := h.Services.Currency.Create(&pgsql.Currency{
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

func (h *Handler) FindCurrencyHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be a number.",
		})
	}

	result, err := h.Services.Currency.FindOneById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "currency not found.",
		})

	}
	return c.JSON(result)
}
