package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, authHttpApi AuthHTTPTransport) {
	authRoutes := router.Group("/auth")
	authRoutes.Post("", authHttpApi.Signup)
	authRoutes.Get("/verify-signup/:token", authHttpApi.VerifySignup)
}
