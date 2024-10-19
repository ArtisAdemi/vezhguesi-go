package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, authHttpApi AuthHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	authRoutes := router.Group("/auth")
	authRoutes.Post("", authHttpApi.Signup)
	authRoutes.Get("/verify-signup/:token", authHttpApi.VerifySignup)
	authRoutes.Post("/login", authHttpApi.Login)
	authRoutes.Put("/update", authMiddleware, authHttpApi.UpdateUser)
}
