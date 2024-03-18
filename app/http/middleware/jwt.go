package middleware

import (
	"errors"
	"time"

	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	JWTTokenAuthed "github.com/golang-jwt/jwt/v4"
)

func JWTProtected() func(*fiber.Ctx) error {
	jwtwareConfig := jwtware.Config{
		SigningKey:     []byte(config.AppCfg().JWTSecretKey),
		ContextKey:     "user", // used in private route
		ErrorHandler:   jwtError,
		SuccessHandler: verifyTokenExpiration,
	}

	return jwtware.New(jwtwareConfig)
}

func verifyTokenExpiration(c *fiber.Ctx) error {
	user := c.Locals("user").(*JWTTokenAuthed.Token)
	claims := user.Claims.(JWTTokenAuthed.MapClaims)
	expires := int64(claims["exp"].(float64))
	if time.Now().Unix() > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": errors.New("token expired"),
		})
	}
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  false,
		"message": err.Error(),
	})
}
