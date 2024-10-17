package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, userHttpApi UserHTTPTransport, authMiddleware func(c *fiber.Ctx) error) {
	userRoutes := router.Group("/users")
	// public routes
	userRoutes.Get("", userHttpApi.GetUsers)

	// Protected with auth middleware
	userRoutes.Get("/:userId", authMiddleware, userHttpApi.GetUserByID)
}
