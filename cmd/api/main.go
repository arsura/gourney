package api

import (
	"fmt"
	"os"

	api "github.com/arsura/gourney/cmd/api/handlers"
	usecase "github.com/arsura/gourney/cmd/usecases"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Application struct {
	Handlers  api.Handlers
	Validator *validator.Validator
	Logger    *zap.SugaredLogger
}

func NewApiApplication(usecases usecase.Usecases, validator *validator.Validator, logger *zap.SugaredLogger) *Application {
	handlers := api.Handlers{
		Currencies: api.NewCurrencyHandler(usecases.Currencies, validator),
	}
	return &Application{
		Handlers:  handlers,
		Validator: validator,
		Logger:    logger,
	}
}

func (app *Application) Start() {
	server := fiber.New()
	server.Use(cors.New())
	app.routes(server)

	port := fmt.Sprintf(":%s", os.Getenv("API_APP_PORT"))
	if err := server.Listen(port); err != nil {
		app.Logger.Errorf("unable to start server: %v", err)
		os.Exit(1)
	}
}
