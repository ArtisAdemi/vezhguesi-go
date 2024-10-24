package reports

import (
	"fmt"
	"vezhguesi/app/entities"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type reportsApi struct {
	db *gorm.DB
	mailDialer *gomail.Dialer
	uiAppUrl string
	logger log.AllLogger
	entitiesApi entities.EntitiesAPI
}

type ReportsAPI interface {
	Create(req *CreateReportRequest) (res *ReportResponse, err error)
	GetReports(req *GetReportsRequest) (res *[]ReportsResponse, err error)
	GetReportByID(req *IDRequest) (res *ReportResponse, err error)
	UpdateReport(req *UpdateReportRequest) (res *ReportResponse, err error)
}

func NewReportsAPI(db *gorm.DB, mailDialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger, entitiesApi entities.EntitiesAPI) ReportsAPI {
	return &reportsApi{db: db, mailDialer: mailDialer, uiAppUrl: uiAppUrl, logger: logger, entitiesApi: entitiesApi}
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
func (s *reportsApi) Create(req *CreateReportRequest) (res *ReportResponse, err error) {
	if req.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if req.StartDate.IsZero() {
		return nil, fmt.Errorf("start date is required")
	}

	if req.EndDate.IsZero() {
		return nil, fmt.Errorf("end date is required")
	}

	report := &Report{
		Subject: req.Subject,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
	}

	result := s.db.Create(&report)
	if result.Error != nil {
		return nil, result.Error
	}

	resp := &ReportResponse{
		Report: *report,
	}

	return resp, nil
}


// @Summary      	Get Reports
// @Description	Validates user id. Gets all reports
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Success			200					{object}	ReportsResponse
// @Router			/api/reports/	[GET]
func (s *reportsApi) GetReports(req *GetReportsRequest) (res *[]ReportsResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	var reports []Report
	result := s.db.Preload("Entities").Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}

	var response []ReportsResponse
	for _, report := range reports {
		response = append(response, ReportsResponse{
			Reports: []Report{report}, // Use 'Reports' and wrap 'report' in a slice
		})
	}

	return &response, nil
}

// @Summary      	Get Report By ID
// @Description	Validates id and user id. Gets report by id
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			id				path		int		true	"Report ID"
// @Success			200					{object}	ReportResponse
// @Router			/api/reports/{id}	[GET]
func (s *reportsApi) GetReportByID(req *IDRequest) (res *ReportResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	if req.ID == 0 {
		return nil, fmt.Errorf("id is required")
	}

	var report Report
	result := s.db.Where("id = ?", req.ID).First(&report)
	if result.Error != nil {
		return nil, result.Error
	}

	resp := &ReportResponse{
		Report: report,
	}

	return resp, nil
}


// @Summary      	Update Report
// @Description	Validates id and user id. Updates report
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			id				path		int		true	"Report ID"
// @Param			UpdateReportRequest	body		UpdateReportRequest	true	"UpdateReportRequest"
// @Success			200					{object}	ReportResponse
// @Router			/api/reports/{id}	[PUT]
func (s *reportsApi) UpdateReport(req *UpdateReportRequest) (res *ReportResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	if req.ID == 0 {
		return nil, fmt.Errorf("id is required")
	}

	var report Report 

	result := s.db.Where("id = ?", req.ID).First(&report)
	if result.Error != nil {
		return nil, fmt.Errorf("report does not exist")
	}

	if req.Title != "" {
		report.Title = req.Title
	}

	if req.Subject != "" {
		report.Subject = req.Subject
	}

	if req.ReportText != "" {
		report.ReportText = req.ReportText
	}

	if req.SourceID != 0 {
		report.SourceID = req.SourceID
	}

	if req.Findings != "" {
		report.Findings = req.Findings
	}

	if len(req.Entities) > 0 {
		var entitiesList []entities.Entity
		for _, entity := range req.Entities {
			requestForEntity := entities.GetEntityRequest{
				Name: entity.Name,
			}
			
			resp, err := s.entitiesApi.GetEntity(&requestForEntity)
			if err != nil {
				 fmt.Println("Entity not found")
			}
			if entity.Name != "" {
				newEntity := entities.CreateEntityRequest{
					Name: entity.Name,
					Type: entity.Type,
				}

				resp, err := s.entitiesApi.Create(&newEntity)
				if err != nil {
					return nil, fmt.Errorf("error creating entity: %v", err)
				}
				createdEntity := entities.Entity{
					ID: resp.ID,
					Name: resp.Name,
					Type: resp.Type,
				}
				entitiesList = append(entitiesList, createdEntity)
			} else {
				entitiesList = append(entitiesList, entities.Entity{
					ID:   resp.ID,
					Name: resp.Name,
					Type: resp.Type,
				})
			}
		}
		report.Entities = entitiesList
	}

	report.Sentiment = req.Sentiment

	result = s.db.Save(&report)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating report: %v", result.Error)
	}

	resp := ReportResponse{
		Report: report,
		UserID: req.UserID,
	}

	return &resp, nil
}
