package api

import "github.com/gofiber/fiber/v2"

func (app *Application) routes(server *fiber.App) {
	server.Get("/currencies/:id", app.Handlers.Currencies.FindCurrencyByIdHandler)
	server.Post("/currencies", app.Handlers.Currencies.CreateCurrencyHandler)
}
