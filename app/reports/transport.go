package reports

import (
	"fmt"
	"strconv"
	"strings"
	"vezhguesi/core/middleware"
	"vezhguesi/helper"

	"github.com/gofiber/fiber/v2"
)

type ReportsHTTPTransport interface {
	Create(c *fiber.Ctx) error
	GetReports(c *fiber.Ctx) error
	GetReportByID(c *fiber.Ctx) error
	UpdateReport(c *fiber.Ctx) error
	GetMyReports(c *fiber.Ctx) error
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

func (s *reportsHttpTransport) GetReports(c *fiber.Ctx) error {
	req := &GetReportsRequest{}
	userId, err := middleware.CtxUserID(c)
	terms := c.Query("terms")
	termsArray := strings.Split(terms, ",")
	if err != nil {
		return helper.HTTPError(c, err, "GetReports.middleware.CtxUserID")
	}
	fmt.Println("terms inside reports transport", terms)
	fmt.Println("termsArray inside reports transport", termsArray)
	req.UserID = userId
	req.Terms = termsArray
	resp, err := s.reportsAPI.GetReports(req)
	if err != nil {
		return helper.HTTPError(c, err, "GetReports.reportsAPI.GetReports")
	}

	return c.JSON(resp)
}

func (s *reportsHttpTransport) GetReportByID(c *fiber.Ctx) error {
	req := &IDRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "GetReportByID.middleware.CtxUserID")
	}
	req.UserID = userId
	reportIdStr := c.Params("id")
	if reportIdStr == "" {
		return helper.HTTPError(c, fmt.Errorf("id is required"), "GetReportByID.c.Params")
	}
	reportId, err := strconv.Atoi(reportIdStr)
	if err != nil {
		return helper.HTTPError(c, err, "GetReportByID.strconv.Atoi")
	}
	req.ID = reportId

	resp, err := s.reportsAPI.GetReportByID(req)
	if err != nil {
		return helper.HTTPError(c, err, "GetReportByID.reportsAPI.GetReportByID")
	}

	return c.JSON(resp)
}

func (s *reportsHttpTransport) UpdateReport(c *fiber.Ctx) error {
	req := &UpdateReportRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "UpdateReport.middleware.CtxUserID")
	}
	req.UserID = userId
	reportIdStr := c.Params("id")
	if reportIdStr == "" {
		return helper.HTTPError(c, fmt.Errorf("id is required"), "UpdateReport.c.Params")
	}
	reportId, err := strconv.Atoi(reportIdStr)
	if err != nil {
		return helper.HTTPError(c, err, "UpdateReport.strconv.Atoi")
	}
	req.ID = reportId

	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "UpdateReport.c.BodyParser")
	}

	resp, err := s.reportsAPI.UpdateReport(req)
	if err != nil {
		return helper.HTTPError(c, err, "UpdateReport.reportsAPI.UpdateReport")
	}
	
	return c.JSON(resp)
	
}

func (s *reportsHttpTransport) GetMyReports(c *fiber.Ctx) error {
	req := &GetReportsRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "GetMyReports.middleware.CtxUserID")
	}
	req.UserID = userId

	resp, err := s.reportsAPI.GetMyReports(req)
	if err != nil {
		return helper.HTTPError(c, err, "GetMyReports.reportsAPI.GetMyReports")
	}

	return c.JSON(resp)
}
