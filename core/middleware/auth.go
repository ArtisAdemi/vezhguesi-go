package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	session "vezhguesi/core/authentication"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func Authentication(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid token")
		}

		// Remove "Bearer " prefix if present
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("userID", claims["userId"])
			return c.Next()
		} else {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}
	}
}

func CtxUserID(c *fiber.Ctx) (int, error) {
	userID, ok := c.Locals("userID").(float64) // JWT claims are often float64
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return int(userID), nil
}

func SessionMiddleware(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        sessionToken := c.Get("Authorization")
        if sessionToken == "" {
            return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
        }

        var session session.Session
        if err := db.Where("session_token = ?", sessionToken).First(&session).Error; err != nil {
            return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
        }

        if session.ExpiresAt.Before(time.Now()) {
            db.Delete(&session)
            return c.Status(http.StatusUnauthorized).SendString("Session expired")
        }

        return c.Next()
    }
}