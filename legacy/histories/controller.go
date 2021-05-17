package histories

import (
	"github.com/gofiber/fiber/v2"
)

func Controller(app *fiber.App) {
	app.Get("/logs", GetHistoryLogsHandler)
}
