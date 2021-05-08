package currencies

import (
	fakedatabase "github.com/arsura/moonbase-service/fake_database"
	"github.com/gofiber/fiber/v2"
)

func GetCurrencyHandler(c *fiber.Ctx) error {
	currencyName := c.Params("name")
	currency := fakedatabase.FindCurrencyByName(currencyName)
	return c.JSON(&fiber.Map{"data": fakedatabase.Currency(currency)})
}
