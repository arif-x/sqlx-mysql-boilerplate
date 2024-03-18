package api

import (
	controllers "github.com/arif-x/sqlx-mysql-boilerplate/app/http/controller/auth"
	"github.com/arif-x/sqlx-mysql-boilerplate/app/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func Auth(a *fiber.App) {
	auth := a.Group("/api/v1/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	need_auth := a.Group("/api/v1/auth", middleware.JWTProtected())
	need_auth.Post("/send-email", controllers.SendEmail)
	need_auth.Get("/verify/:token", controllers.Verify)
}
