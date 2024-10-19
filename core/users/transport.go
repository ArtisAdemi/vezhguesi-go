package users

import (
	"strconv"
	"vezhguesi/core/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHTTPTransport interface {
	GetUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetUserData(c *fiber.Ctx) error
}

type userHttpTransport struct {
	userAPI UserAPI
}

func NewUserHTTPTransport(userAPI UserAPI) UserHTTPTransport {
	return &userHttpTransport{userAPI: userAPI}
}


func (s *userHttpTransport) GetUsers(c *fiber.Ctx) error {
	req := &FindRequest{}
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := s.userAPI.GetUsers(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (s *userHttpTransport) GetUserByID(c *fiber.Ctx) error {
	req := &FindUserByID{}
	userIdParamStr := c.Params("userId")
	userId, err := strconv.Atoi(userIdParamStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.UserID = userId
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := s.userAPI.GetUserByID(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (s *userHttpTransport) GetUserData(c *fiber.Ctx) error {
	req := &FindUserByID{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	req.UserID = userId

	resp, err := s.userAPI.GetUserData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
