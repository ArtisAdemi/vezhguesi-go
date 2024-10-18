package entities

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EntitiesHTTPTransport interface {
	Create(c *fiber.Ctx) error
	GetEntity(c *fiber.Ctx) error
}

type entitiesHttpTransport struct {
	entitiesAPI EntitiesAPI
}

func NewEntitiesHTTPTransport(entitiesAPI EntitiesAPI) EntitiesHTTPTransport {
	return &entitiesHttpTransport{entitiesAPI: entitiesAPI}
}

func (s *entitiesHttpTransport) Create(c *fiber.Ctx) error {
	req := &CreateEntityRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	res, err := s.entitiesAPI.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}


func (s *entitiesHttpTransport) GetEntity(c *fiber.Ctx) error {
	req := &GetEntityRequest{}
	idStr := c.Params("id")
	idInt, _ := strconv.ParseUint(idStr, 10, 64)
	req.ID = uint(idInt)

	name := c.Query("name")
	req.Name = name

	res, err := s.entitiesAPI.GetEntity(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

