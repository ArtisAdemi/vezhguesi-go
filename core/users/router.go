package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, userHttpApi UserHTTPTransport) {
	router.Get("/users", userHttpApi.GetUsers)
	router.Post("/users", userHttpApi.Signup)
}
