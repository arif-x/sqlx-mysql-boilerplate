package middleware

import (
	"github.com/gofiber/fiber/v2"
	JWTTokenAuthed "github.com/golang-jwt/jwt/v4"
)

func Email() func(*fiber.Ctx) error {
	middleware := func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTTokenAuthed.Token)
		claims := user.Claims.(JWTTokenAuthed.MapClaims)
		email_verified_at := claims["email_verified_at"]

		if email_verified_at != nil {
			return c.Next()
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Please verify your email first!",
			})
		}
	}
	return middleware
}
