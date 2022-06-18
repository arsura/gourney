package main

import (
	api "github.com/arsura/gourney/cmd/api"
	usecase "github.com/arsura/gourney/cmd/usecases"
	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/logger"
	repository "github.com/arsura/gourney/pkg/repositories"
	"github.com/arsura/gourney/pkg/validator"
)

func main() {

	logger := logger.NewLogger()
	validator := validator.NewValidator()
	config := config.NewConfig()
	mongoClient := adapter.NewMongoDBClient(config)

	postRepository := repository.NewPostRepository(mongoClient, logger, config)
	repositories := repository.Repositories{
		Posts: postRepository,
	}

	postUsecase := usecase.NewPostUsecase(&repositories, logger)
	usecases := usecase.Usecases{
		Post: postUsecase,
	}

	if isApiEnable := config.APIServer.IsEnable; isApiEnable {
		apiApp := api.NewApiApplication(usecases, validator, logger, config)
		apiApp.Start()
	}

	// if isRabbitMqEnable ?
	// if isKafkaEnable ?
	// if bla bla app enable

}
