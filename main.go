package main

import (
	"log"

	"github.com/arsura/moonbase-service/currencies"
	"github.com/arsura/moonbase-service/histories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	currencies.Controller(app)
	histories.Controller(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
