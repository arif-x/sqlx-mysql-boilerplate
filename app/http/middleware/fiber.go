package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			AllowOrigins: "https://gofiber.io, https://gofiber.net",
			AllowHeaders: "Origin, Content-Type, Accept",
		}),
		// Add simple logger.
		logger.New(),
	)
}
