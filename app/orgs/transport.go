package orgs

import (
	"vezhguesi/core/middleware"
	"vezhguesi/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OrgHTTPTransport interface {
	Add(c *fiber.Ctx) error
}

type orgHttpTransport struct {
	orgApi OrgAPI
	logger log.AllLogger
}

func NewOrgHTTPTransport(orgApi OrgAPI, logger log.AllLogger) OrgHTTPTransport {
	return &orgHttpTransport{
		orgApi: orgApi,
		logger: logger,
	}
}

func (s *orgHttpTransport) Add(c *fiber.Ctx) error {
	req := &AddOrgRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.CtxUserID")
	}
	req.UserID = userId
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.BodyParser")
	}

	resp, err := s.orgApi.Add(req)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.Add")
	}

	return c.JSON(resp)
}
