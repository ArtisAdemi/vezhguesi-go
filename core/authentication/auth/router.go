package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, authHttpApi AuthHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	authRoutes := router.Group("/auth")
	// Public routes
	authRoutes.Post("", authHttpApi.Signup)
	authRoutes.Get("/verify-signup/:token", authHttpApi.VerifySignup)
	authRoutes.Post("/login", authHttpApi.Login)
	authRoutes.Post("/forgot-password", authHttpApi.ForgotPassword)
	authRoutes.Put("/reset-password/:token", authHttpApi.ResetPassword)
	// Protected routes
	authRoutes.Put("/update", authMiddleware, authHttpApi.UpdateUser)
}
