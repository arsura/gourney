package main

import (
	"context"
	"fmt"

	api "github.com/arsura/gourney/cmd/api"
	usecase "github.com/arsura/gourney/cmd/usecases"
	"github.com/arsura/gourney/config"
	"github.com/arsura/gourney/pkg/logger"
	repository "github.com/arsura/gourney/pkg/repositories"
	"github.com/arsura/gourney/pkg/validator"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	logger := logger.NewLogger()
	validator := validator.NewValidator()
	config := config.NewConfig()

	pool, err := pgxpool.Connect(context.Background(), config.PostgresURI)
	if err != nil {
		panic(fmt.Errorf("unable to connect database: %w", err))
	}
	defer pool.Close()

	repositories := repository.Repositories{
		Currencies: repository.NewCurrencyRepo(pool),
	}

	usercases := usecase.Usecases{
		Currencies: usecase.NewCurrencyUsecase(repositories.Currencies, logger),
	}

	if isApiEnable := config.APIServer.IsEnable; isApiEnable {
		apiApp := api.NewApiApplication(usercases, validator, logger, config)
		apiApp.Start()
	}

	// if isRabbitMqEnable ?
	// if isKafkaEnable ?
	// if bla bla app enable

}
