package main

import (
	"context"
	"fmt"
	"os"

	handler "github.com/arsura/moonbase-service/cmd/handlers"
	service "github.com/arsura/moonbase-service/cmd/services"
	"github.com/arsura/moonbase-service/pkg/logger"
	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	"github.com/arsura/moonbase-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	handler *handler.Handlers
}

func main() {
	server := fiber.New()
	server.Use(cors.New())

	logger := logger.InitLogger()

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Errorf("Unable to connect database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	validate, trans := validator.InitValidator()
	app := &application{
		handler: &handler.Handlers{
			Currencies: &handler.CurrencyHandler{
				Validator: &validator.Validator{
					Validate: validate,
					Trans:    trans,
				},
				CurrencyService: &service.CurrencyService{
					Logger:       logger,
					CurrencyRepo: &pgsql.CurrencyRepo{Conn: pool},
				},
			},
		},
	}
	app.routes(server)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := server.Listen(port); err != nil {
		logger.Errorf("Unable to start server: %v", err)
		os.Exit(1)
	}
}
