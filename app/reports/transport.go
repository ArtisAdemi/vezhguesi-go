package reports

import (
	"vezhguesi/core/middleware"
	"vezhguesi/helper"

	"github.com/gofiber/fiber/v2"
)

type ReportsHTTPTransport interface {
	Create(c *fiber.Ctx) error
}

type reportsHttpTransport struct {
	reportsAPI ReportsAPI
}

func NewReportsHTTPTransport(reportsAPI ReportsAPI) ReportsHTTPTransport {
	return &reportsHttpTransport{reportsAPI: reportsAPI}
}

func (s *reportsHttpTransport) Create(c *fiber.Ctx) error {
	req := &CreateReportRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "CreateReport.middleware.CtxUserID")
	}

	req.UserID = userId

	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "CreateReport.c.BodyParser")
	}

	resp, err := s.reportsAPI.Create(req)
	if err != nil {
		return helper.HTTPError(c, err, "CreateReport.reportsAPI.Create")
	}

	return c.JSON(resp)
}