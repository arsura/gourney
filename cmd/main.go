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
	logger  *logger.Logger
	handler *handler.Handler
}

func main() {
	server := fiber.New()
	server.Use(cors.New())

	logger := logger.InitLog()
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error.Printf("Unable to connect database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	service := &service.Service{
		Logger: logger,
		Pg: &pgsql.Repositories{
			Currencies: &pgsql.DB{Conn: pool},
		},
	}

	validate, trans := validator.InitValidate()
	app := &application{
		logger: logger,
		handler: &handler.Handler{
			Service: service,
			Validator: &validator.Validator{
				Validate: validate,
				Trans:    trans,
			},
		},
	}
	app.routes(server)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := server.Listen(port); err != nil {
		logger.Error.Printf("Unable to start server: %v\n", err)
		os.Exit(1)
	}
}
