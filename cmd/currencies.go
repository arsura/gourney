package main

import (
	"strconv"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
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

	_, err := app.pg.Currencies.Insert(&pgsql.Currency{
		Name:       newCurrency.Name,
		Amount:     newCurrency.Amount,
		Total:      newCurrency.Total,
		RiseRate:   newCurrency.RiseRate,
		RiseFactor: newCurrency.RiseFactor,
	})

	if err != nil {
		app.errorLog.Printf("Failed to create currency: %v\n", err)
		return c.Status(400).JSON(&fiber.Map{
			"error": "Failed to create currency.",
		})
	}

	return c.SendStatus(201)
}

func (app *application) findCurrency(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		app.errorLog.Printf("Failed to create currency: %v\n", err)
		return c.Status(400).JSON(&fiber.Map{
			"error": "Invalid currency param, it should be a number.",
		})
	}

	result, err := app.pg.Currencies.Get(int64(id))
	if err != nil {
		app.errorLog.Printf("Failed to find currency: %v\n", err)
		return c.Status(404).JSON(&fiber.Map{
			"error": "Currency not found.",
		})

	}
	return c.JSON(result)
}
