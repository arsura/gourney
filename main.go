package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	currencies *pgsql.CurrencyModel
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	appEnv := &application{
		currencies: &pgsql.CurrencyModel{Pool: pool},
	}

	fmt.Println(appEnv)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
