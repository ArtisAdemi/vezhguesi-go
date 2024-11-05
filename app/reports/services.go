package reports

import (
	"fmt"
	"strings"
	"vezhguesi/app/entities"
	server "vezhguesi/sentiment-communication"

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
	sentiment server.ServerAPI
}

type ReportsAPI interface {
	Create(req *CreateReportRequest) (res *ReportResponse, err error)
	GetReports(req *GetReportsRequest) (res *GetReportsResponse, err error)
	GetReportByID(req *IDRequest) (res *ReportResponse, err error)
	UpdateReport(req *UpdateReportRequest) (res *ReportResponse, err error)
	GetMyReports(req *GetReportsRequest) (res *GetMyReportsResponse, err error)
}

func NewReportsAPI(db *gorm.DB, mailDialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger, entitiesApi entities.EntitiesAPI, serverApi server.ServerAPI) ReportsAPI {
	return &reportsApi{db: db, mailDialer: mailDialer, uiAppUrl: uiAppUrl, logger: logger, entitiesApi: entitiesApi, sentiment: serverApi}
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
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	if req.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if req.StartDate.IsZero() {
		return nil, fmt.Errorf("start date is required")
	}

	if req.EndDate.IsZero() {
		return nil, fmt.Errorf("end date is required")
	}

	subjectList := strings.Split(req.Subject, ",")

	report := &Report{
		Subject: req.Subject,
		UserID: req.UserID,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
	}

	articles, err := s.sentiment.FetchArticles()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %v", err)
	}

	var articlesList []Articles
	for _, article := range articles {
		for _, subject := range subjectList {
			if strings.Contains(article.Content, subject) {
				articlesList = append(articlesList, Articles{
					ID: article.ID,
					Title: article.Title,
					Content: article.Content,
				})
			}
		}
	}

	var articleIds []int
	for _, article := range articlesList {
		articleIds = append(articleIds, article.ID)
	}

	_, err = s.sentiment.AnalyzeArticles(&articleIds)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze articles: %v", err)
	}

	result := s.db.Create(&report)
	if result.Error != nil {
		return nil, result.Error
	}


	resp := &ReportResponse{
		Report: *report,
		Articles: articlesList,
	}

	return resp, nil
}


// @Summary      	Get Reports
// @Description	Validates user id. Gets all reports
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			terms			query		string	true	"terms"
// @Success			200					{object}	GetReportsResponse
// @Router			/api/reports/	[GET]
func (s *reportsApi) GetReports(req *GetReportsRequest) (res *GetReportsResponse, err error) {
	// Call the GetAnalyzes function
	analyzeResponse, err := s.sentiment.GetAnalyzes(req.Terms)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analyzed reports: %v", err)
	}

	// Transform the analyzeResponse into GetReportsResponse
	analysesResponse := &GetReportsResponse{
		Analyses:      make([]Analysis, 0),
		TotalArticles: analyzeResponse.Results.TotalArticles,
	}

	for _, article := range analyzeResponse.Results.Articles {
		// Create a new Analysis object
		analysis := Analysis{
			AnalysisResults: AnalysisResults{
				Entities: make([]AnalysisEntity, 0),
				Topics:   make([]AnalysisTopic, 0),
			},
			ArticleMetadata: ArticleMetadata{
				ArticleSummary: article.ArticleSummary,
				ID:            article.ArticleID,
				PublishedDate: article.PublishedDate,
				ScrapedAt:     article.ScrapedAt,
				Title:         article.Title,
				URLID:         article.URLID,
				URL:          article.URL,
			},
		}

		// Convert map of entities to slice of AnalysisEntity
		for _, entity := range article.Entities {
			analysisEntity := AnalysisEntity{
				Name:           entity.Name,
				RelatedTopics:  entity.RelatedTopics,
				SentimentLabel: entity.SentimentLabel,
				SentimentScore: entity.SentimentScores[0],
			}
			analysis.AnalysisResults.Entities = append(analysis.AnalysisResults.Entities, analysisEntity)
		}

		// Add the analysis to the response
		analysesResponse.Analyses = append(analysesResponse.Analyses, analysis)
	}

	return analysesResponse, nil
}

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

// @Summary      	Get My Reports
// @Description	Validates user id. Gets all reports made by the user
// @Tags			Reports
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Success			200					{object}	GetMyReportsResponse
// @Router			/api/reports/my-reports	[GET]
func (s *reportsApi) GetMyReports(req *GetReportsRequest) (res *GetMyReportsResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	var reports []Report
	result := s.db.Where("user_id = ?", req.UserID).Find(&reports)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching reports: %v", result.Error)
	}

	if len(reports) == 0 {
		return nil, fmt.Errorf("no reports found")
	}

	var terms []string
	for _, report := range reports {
		terms = append(terms, report.Subject)
	}

	fmt.Println("______________________")
	fmt.Println("terms: ", terms)

	response, err := s.sentiment.GetAnalyzes(terms)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analyzed reports: %v", err)
	}

	// Group analyses by entity
	entityMap := make(map[string]*EntityAnalysis)

	for _, article := range response.Results.Articles {
		analysis := createAnalysisFromArticle(article)

		// Add this analysis to each entity it mentions
		for _, entity := range article.Entities {
			if _, exists := entityMap[entity.Name]; !exists {
				entityMap[entity.Name] = &EntityAnalysis{
					EntityName: entity.Name,
					Analyses:   []Analysis{},
				}
			}
			entityMap[entity.Name].Analyses = append(entityMap[entity.Name].Analyses, analysis)
		}
	}

	// Convert map to slice
	var entities []EntityAnalysis
	for _, entityAnalysis := range entityMap {
		entityAnalysis.TotalArticles = len(entityAnalysis.Analyses)
		entities = append(entities, *entityAnalysis)
	}

	return &GetMyReportsResponse{
		Entities: entities,
	}, nil
}

// Helper function to create Analysis from ArticleData
func createAnalysisFromArticle(article server.ArticleData) Analysis {
	var analysisResults AnalysisResults
	
	// Convert entities
	for _, entity := range article.Entities {
		analysisEntity := AnalysisEntity{
			Name:           entity.Name,
			RelatedTopics:  entity.RelatedTopics,
			SentimentLabel: entity.SentimentLabel,
			SentimentScore: entity.SentimentScores[0],
		}
		analysisResults.Entities = append(analysisResults.Entities, analysisEntity)
	}

	// Convert topics
	for _, topic := range article.Topics {
		analysisTopic := AnalysisTopic{
			Name:           topic.Name,
			SentimentLabel: topic.SentimentLabel,
			SentimentScore: topic.SentimentScore,
		}
		analysisResults.Topics = append(analysisResults.Topics, analysisTopic)
	}

	return Analysis{
		AnalysisResults: analysisResults,
		ArticleMetadata: ArticleMetadata{
			ArticleSummary: article.ArticleSummary,
			ID:            article.ArticleID,
			PublishedDate: article.PublishedDate,
			ScrapedAt:     article.ScrapedAt,
			Title:         article.Title,
			URLID:         article.URLID,
			URL:           article.URL,
		},
	}
}
