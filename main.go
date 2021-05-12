package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arsura/moonbase-service/currencies"
	"github.com/arsura/moonbase-service/histories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	currencies.Controller(app)
	histories.Controller(app)

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
