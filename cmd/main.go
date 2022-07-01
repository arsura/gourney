package main

import (
	"context"

	"github.com/arsura/gourney/cmd/api"
	"github.com/arsura/gourney/cmd/logsworker"
	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/logger"
	repository "github.com/arsura/gourney/pkg/repositories"
	service "github.com/arsura/gourney/pkg/services"
	usecase "github.com/arsura/gourney/pkg/usecases"
	"github.com/arsura/gourney/pkg/validator"
)

func main() {

	var (
		logger             = logger.NewLogger()
		validator          = validator.NewValidator()
		config             = config.NewConfig(logger)
		mongoClient        = adapter.NewMongoDBClient(logger, config)
		rabbitMQConnection = adapter.NewRabbitMQConnection(logger, config)
	)

	defer mongoClient.Disconnect(context.Background())
	defer rabbitMQConnection.Connection.Close()
	defer rabbitMQConnection.Channel.Close()

	var (
		postRepository    = repository.NewPostRepository(mongoClient, logger, config)
		postLogRepository = repository.NewPostLogRepository(mongoClient, logger, config)
		repositories      = &repository.Repository{
			Post:    postRepository,
			PostLog: postLogRepository,
		}
	)

	var (
		logService = service.NewLogService(rabbitMQConnection, logger, config)
		services   = &service.Services{
			Log: logService,
		}
	)

	var (
		postUsecase    = usecase.NewPostUsecase(repositories, services, logger)
		postLogUsecase = usecase.NewPostLogUsecase(repositories, services, logger)
		usecases       = &usecase.Usecase{
			Post:    postUsecase,
			PostLog: postLogUsecase,
		}
	)

	if isApiEnable := config.APIService.IsEnable; isApiEnable {
		apiApp := api.NewApiApplication(usecases, validator, logger, config)
		apiApp.Start()
	}

	if isLogsWorkerEnable := config.LogsWorkerService.IsEnable; isLogsWorkerEnable {
		workerApp := logsworker.NewWorkerApplication(rabbitMQConnection, usecases, logger, config)
		workerApp.Start()
	}
}
