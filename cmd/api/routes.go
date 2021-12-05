package api

import "github.com/gofiber/fiber/v2"

func (app *Application) routes(server *fiber.App) {
	server.Get("/currencies/:id", app.Handler.Currencies.FindCurrencyByIdHandler)
	server.Post("/currencies", app.Handler.Currencies.CreateCurrencyHandler)
}
