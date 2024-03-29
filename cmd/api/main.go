package api

import (
	"fmt"

	api "github.com/arsura/gourney/cmd/api/handlers"
	middleware "github.com/arsura/gourney/cmd/api/middlewares"
	"github.com/arsura/gourney/config"
	usecase "github.com/arsura/gourney/pkg/usecases"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

type Application struct {
	Handlers  *api.Handler
	Validator *validator.Validator
	Logger    *zap.SugaredLogger
	Config    *config.Config
}

func NewApiApplication(usecases *usecase.UseCase, validator *validator.Validator, logger *zap.SugaredLogger, config *config.Config) *Application {
	handlers := &api.Handler{
		Post: api.NewPostHandler(usecases.Post, validator, logger),
	}
	return &Application{handlers, validator, logger, config}
}

func (app *Application) Start() {
	server := fiber.New()
	server.Use(cors.New())
	server.Use(requestid.New())
	server.Use(middleware.RequestLogging())
	app.routes(server)

	port := fmt.Sprintf(":%s", app.Config.APIService.Port)
	if err := server.Listen(port); err != nil {
		app.Logger.With("error", err).Panic("unable to start server")
	}
}
