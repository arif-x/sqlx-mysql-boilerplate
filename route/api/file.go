package api

import "github.com/gofiber/fiber/v2"

func FileRoutes(a *fiber.App) {
	a.Static("/upload/post", "./upload/post")
}
