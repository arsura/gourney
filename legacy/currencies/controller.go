package currencies

import (
	"github.com/gofiber/fiber/v2"
)

func Controller(app *fiber.App) {
	app.Get("/currencies/:name", GetCurrencyHandler)
	app.Get("/amounts", CalcAmountsHandler)
	app.Get("/rates", GetExchangeRateHandler)
	app.Post("/purchases", PurchaseHandler)
}
