package main

import "github.com/gofiber/fiber/v2"

func (app *application) routes(server *fiber.App) {
	server.Get("/currencies/:id", app.handler.FindCurrencyHandler)
	server.Post("/currencies", app.handler.CreateCurrencyHandler)
}
