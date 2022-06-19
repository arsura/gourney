package api

import (
	"fmt"
	"os"

	api "github.com/arsura/gourney/cmd/api/handlers"
	middleware "github.com/arsura/gourney/cmd/api/middlewares"
	config "github.com/arsura/gourney/configs"
	usecase "github.com/arsura/gourney/pkg/usecases"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

type Application struct {
	Handlers  api.Handlers
	Validator *validator.Validator
	Logger    *zap.SugaredLogger
	Config    *config.Config
}

func NewApiApplication(usecases usecase.Usecases, validator *validator.Validator, logger *zap.SugaredLogger, config *config.Config) *Application {
	handlers := api.Handlers{
		Post: api.NewPostHandler(usecases.Post, validator, logger),
	}
	return &Application{
		Handlers:  handlers,
		Validator: validator,
		Logger:    logger,
		Config:    config,
	}
}

func (app *Application) Start() {
	server := fiber.New()
	server.Use(cors.New())
	server.Use(requestid.New())
	server.Use(middleware.RequestLogging())
	app.routes(server)

	port := fmt.Sprintf(":%s", app.Config.APIServer.Port)
	if err := server.Listen(port); err != nil {
		app.Logger.Errorf("unable to start server: %v", err)
		os.Exit(1)
	}
}
