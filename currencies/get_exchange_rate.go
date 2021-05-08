package currencies

import (
	fakedatabase "github.com/arsura/moonbase-service/fake_database"
	"github.com/gofiber/fiber/v2"
)

func GetExchangeRateHandler(c *fiber.Ctx) error {
	query := new(CalcAmountQuery)
	if err := c.QueryParser(query); err != nil {
		return err
	}
	sourceCurrency := query.Source
	targetCurrency := query.Target
	rate := fakedatabase.FindExchangeRate(sourceCurrency, targetCurrency)
	return c.JSON(&fiber.Map{"data": rate})
}
