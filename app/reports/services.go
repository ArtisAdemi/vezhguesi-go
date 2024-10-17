package reports

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type reportsApi struct {
	db *gorm.DB
	mailDialer *gomail.Dialer
	uiAppUrl string
	logger log.AllLogger
}

type ReportsAPI interface {
	Create(req *CreateReportRequest) (res *ReportResponse, err error)
}

func NewReportsAPI(db *gorm.DB, mailDialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger) ReportsAPI {
	return &reportsApi{db: db, mailDialer: mailDialer, uiAppUrl: uiAppUrl, logger: logger}
}

// @Summary      	Create Report
// @Description	Validates subject, start date, end date. Creates a new report.
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			CreateReportRequest	body		CreateReportRequest	true	"CreateReportRequest"
// @Success			200					{object}	ReportResponse
// @Router			/api/reports/	[POST]
func (r *reportsApi) Create(req *CreateReportRequest) (res *ReportResponse, err error) {
	if req.Subject == "" {
		return nil, errors.New("subject is required")
	}

	if req.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}

	if req.EndDate.IsZero() {
		return nil, errors.New("end date is required")
	}

	report := Report{
		Subject: req.Subject,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
	}

	resp := &ReportResponse{
		Report: report,
	}

	return resp, nil
}

