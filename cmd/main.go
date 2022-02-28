package main

import (
	"context"
	"os"
	"strconv"

	api "github.com/arsura/gourney/cmd/api"
	usecase "github.com/arsura/gourney/cmd/usecases"
	"github.com/arsura/gourney/pkg/logger"
	repository "github.com/arsura/gourney/pkg/repositories"
	"github.com/arsura/gourney/pkg/validator"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	logger := logger.NewLogger()
	validator := validator.NewValidator()

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URI"))
	if err != nil {
		logger.Errorf("unable to connect database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	repositories := repository.Repositories{
		Currencies: repository.NewCurrencyRepo(pool),
	}

	usercases := usecase.Usecases{
		Currencies: usecase.NewCurrencyUsecase(repositories.Currencies, logger),
	}

	if isApiEnable, err := strconv.ParseBool(os.Getenv("IS_API_ENABLE")); err == nil && isApiEnable {
		apiApp := api.NewApiApplication(usercases, validator, logger)
		apiApp.Start()
	}

	// if isRabbitMqEnable ?
	// if isKafkaEnable ?
	// if bla bla app enable

}
