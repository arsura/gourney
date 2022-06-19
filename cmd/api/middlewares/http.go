package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RequestLogging() func(*fiber.Ctx) error {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${locals:requestid} | ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	})
}
