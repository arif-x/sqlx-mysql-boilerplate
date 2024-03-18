package middleware

import (
	"github.com/gofiber/fiber/v2"
	JWTTokenAuthed "github.com/golang-jwt/jwt/v4"
)

func Permission(Permission string) func(*fiber.Ctx) error {
	middleware := func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTTokenAuthed.Token)
		claims := user.Claims.(JWTTokenAuthed.MapClaims)
		permissions := claims["permission"]

		permissionsSlice, ok := permissions.([]interface{})
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "You dont have permission!",
			})
		}

		var permissionsArray []string
		for _, v := range permissionsSlice {
			if s, ok := v.(string); ok {
				permissionsArray = append(permissionsArray, s)
			}
		}

		found := false
		for _, check := range permissionsArray {
			if check == Permission {
				found = true
				break
			}
		}

		if found {
			return c.Next()
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "You dont have permission!",
			})
		}
	}

	return middleware
}
