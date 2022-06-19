package main

import (
	"context"

	api "github.com/arsura/gourney/cmd/api"
	config "github.com/arsura/gourney/configs"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/logger"
	repository "github.com/arsura/gourney/pkg/repositories"
	usecase "github.com/arsura/gourney/pkg/usecases"
	"github.com/arsura/gourney/pkg/validator"
)

func main() {

	var (
		logger             = logger.NewLogger()
		validator          = validator.NewValidator()
		config             = config.NewConfig(logger)
		mongoClient        = adapter.NewMongoDBClient(logger, config)
		rabbitMqConnection = adapter.NewRabbitMQConnection(logger, config)
	)

	var (
		postRepository = repository.NewPostRepository(mongoClient, logger, config)

		repositories = repository.Repositories{
			Posts: postRepository,
		}
	)

	var (
		postUsecase = usecase.NewPostUsecase(&repositories, rabbitMqConnection, logger)

		usecases = usecase.Usecases{
			Post: postUsecase,
		}
	)

	if isApiEnable := config.APIServer.IsEnable; isApiEnable {
		apiApp := api.NewApiApplication(usecases, validator, logger, config)
		apiApp.Start()
	}

	// if isRabbitMqEnable ?
	// if isKafkaEnable ?
	// if bla bla app enable

	defer mongoClient.Disconnect(context.Background())
	defer rabbitMqConnection.Connection.Close()
	defer rabbitMqConnection.Channel.Close()
}
