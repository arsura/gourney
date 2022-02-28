package main

import (
	"context"
	"os"
	"strconv"

	api "github.com/arsura/gourney/cmd/api"
	"github.com/arsura/gourney/pkg/logger"
	"github.com/arsura/gourney/pkg/validator"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	logger := logger.InitLogger()
	validatr, trans := validator.InitValidator()

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Errorf("Unable to connect database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if isApiEnable, err := strconv.ParseBool(os.Getenv("IS_API_ENABLE")); err == nil && isApiEnable {
		api := &api.Application{
			Validator: &validator.Validator{
				Validate: validatr,
				Trans:    trans,
			},
			Logger: logger,
			DbConn: pool,
		}
		api.Start()
	}

	// if isRabbitMqEnable ?
	// if isKafkaEnable ?
	// if bla bla app enable

}
