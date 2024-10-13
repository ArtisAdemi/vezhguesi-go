package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, userHttpApi UserHTTPTransport) {
	userRoutes := router.Group("/users")
	userRoutes.Get("", userHttpApi.GetUsers)
	userRoutes.Post("", userHttpApi.Signup)
	userRoutes.Get("/verify-signup/:token", userHttpApi.VerifySignup)
}
