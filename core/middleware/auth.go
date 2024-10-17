package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type HTTPError struct {
	Message string `json:"message"`
}

//  Protects protected routes
func Authentication(key string) func(*fiber.Ctx) error{
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(key),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(HTTPError{Message: ErrUnauthorized.Error()})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(HTTPError{Message: err.Error()})
	}
}