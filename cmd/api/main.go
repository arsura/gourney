package api

import (
	"fmt"
	"os"

	handler "github.com/arsura/gourney/cmd/api/handlers"
	usecase "github.com/arsura/gourney/cmd/usecases"
	repo "github.com/arsura/gourney/pkg/repositories"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Application struct {
	Handler   *handler.Handlers
	Validator *validator.Validator
	Logger    *zap.SugaredLogger
	Conn      repo.DbConn
}

func (app *Application) Start() {
	server := fiber.New()
	server.Use(cors.New())

	api := &Application{
		Handler: &handler.Handlers{
			Currencies: &handler.CurrencyHandler{
				Validator: app.Validator,
				CurrencyUsecase: &usecase.CurrencyUsecase{
					Logger:       app.Logger,
					CurrencyRepo: &repo.CurrencyRepo{Conn: app.Conn},
				},
			},
		},
	}
	api.routes(server)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := server.Listen(port); err != nil {
		app.Logger.Errorf("Unable to start server: %v", err)
		os.Exit(1)
	}
}
