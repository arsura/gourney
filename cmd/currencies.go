package main

import (
	"github.com/arsura/moonbase-service/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type createReqBody struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Total      float64 `json:"total"`
	RiseRate   float64 `json:"riseRate"`
	RiseFactor float64 `json:"riseFactor"`
}

func (app *application) createCurrency(c *fiber.Ctx) error {
	newCurrency := new(createReqBody)
	if err := c.BodyParser(newCurrency); err != nil {
		return err
	}

	_, err := app.currencies.Insert(&models.Currency{
		Name:       newCurrency.Name,
		Amount:     newCurrency.Amount,
		Total:      newCurrency.Total,
		RiseRate:   newCurrency.RiseRate,
		RiseFactor: newCurrency.RiseFactor,
	})

	if err != nil {
		app.errorLog.Printf("Failed to create currency: %v\n", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(&fiber.Map{"ping": newCurrency})
}
