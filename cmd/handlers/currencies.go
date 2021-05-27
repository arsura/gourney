package handler

import (
	"strconv"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	"github.com/gofiber/fiber/v2"
)

type CreateReqBody struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Total      float64 `json:"total"`
	RiseRate   float64 `json:"riseRate"`
	RiseFactor float64 `json:"riseFactor"`
}

func (h *Handler) CreateCurrencyHandler(c *fiber.Ctx) error {
	currency := new(CreateReqBody)
	if err := c.BodyParser(currency); err != nil {
		return err
	}

	_, err := h.Service.CreateCurrency(&pgsql.Currency{
		Name:       currency.Name,
		Amount:     currency.Amount,
		Total:      currency.Total,
		RiseRate:   currency.RiseRate,
		RiseFactor: currency.RiseFactor,
	})

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create currency.",
		})
	}

	return c.SendStatus(201)
}

func (h *Handler) FindCurrencyHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid param, it should be a number.",
		})
	}

	result, err := h.Service.FindCurrency(int64(id))
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"error": "Currency not found.",
		})

	}
	return c.JSON(result)
}
