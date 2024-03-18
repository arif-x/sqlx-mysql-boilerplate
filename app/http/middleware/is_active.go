package middleware

import (
	"github.com/gofiber/fiber/v2"
	JWTTokenAuthed "github.com/golang-jwt/jwt/v4"
)

func IsActive() func(*fiber.Ctx) error {
	middleware := func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTTokenAuthed.Token)
		claims := user.Claims.(JWTTokenAuthed.MapClaims)
		is_active := claims["is_active"]

		if is_active == true {
			return c.Next()
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "You are not an active user!",
			})
		}
	}
	return middleware
}
