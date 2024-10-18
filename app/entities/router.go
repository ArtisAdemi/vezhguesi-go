package entities

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, transport EntitiesHTTPTransport) {
	router.Post("/entities", transport.Create)
	router.Get("/entities/:id", transport.GetEntity)
}