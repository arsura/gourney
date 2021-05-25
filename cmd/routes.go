package main

import "github.com/gofiber/fiber/v2"

func (app *application) routes(server *fiber.App) {
	server.Get("/currencies/:id", app.findCurrency)
	server.Post("/currencies", app.createCurrency)
}
