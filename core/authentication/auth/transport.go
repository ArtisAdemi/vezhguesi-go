package auth

import "github.com/gofiber/fiber/v2"

type AuthHTTPTransport interface {
	Signup(c *fiber.Ctx) error
	VerifySignup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authHttpTransport struct {
	authAPI AuthApi
}

func NewAuthHTTPTransport(authAPI AuthApi) AuthHTTPTransport {
	return &authHttpTransport{authAPI: authAPI}
}

func (s *authHttpTransport) Signup(c *fiber.Ctx) error {
	req := &SignupRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := s.authAPI.Signup(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (s *authHttpTransport) VerifySignup(c *fiber.Ctx) error {
	req := &SignupVerifyRequest{}

	req.Token = c.Params("token")
	resp, err := s.authAPI.VerifySignup(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (s *authHttpTransport) Login(c *fiber.Ctx) error {
	req := &LoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	resp, err := s.authAPI.Login(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

