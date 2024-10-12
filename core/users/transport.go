package users

import "github.com/gofiber/fiber/v2"

type UserHTTPTransport interface {
	GetUsers(c *fiber.Ctx) error
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
