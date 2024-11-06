package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"vezhguesi/app/entities"
	entity_reportsvc "vezhguesi/app/entity_reports"
	"vezhguesi/helper"
	server "vezhguesi/sentiment-communication"

	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/sashabaranov/go-openai"
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
				SentimentScore: entity.SentimentScore,
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


	response, err := s.sentiment.GetAnalyzes(terms)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analyzed reports: %v", err)
	}

	// Group analyses by entity and track sentiment scores
	entityMap := make(map[string]*EntityAnalysis)
	entitySentiments := make(map[string][]float32)
	
	// Store all articles for later summary generation
	articles := response.Results.Articles
	
	var articlesList []string
	for _, article := range articles {
		analysis := createAnalysisFromArticle(article)

		for entityName, entity := range article.Entities {
			articlesList = append(articlesList, article.URL)
			
			if _, exists := entityMap[entityName]; !exists {
				entityMap[entityName] = &EntityAnalysis{
						EntityName: entityName,
						Analyses:   []Analysis{},
				}
				entitySentiments[entityName] = []float32{}
			}
			entityMap[entityName].Analyses = append(entityMap[entityName].Analyses, analysis)
			
			// Use the sentiment score directly from the entity
			entitySentiments[entityName] = append(entitySentiments[entityName], entity.SentimentScore)
		}
	}


	var entitiesReportsResponse []EntityReport

	// Convert map to slice, calculate averages, and generate summaries
	for entityName, entityAnalysis := range entityMap {
		// Calculate sentiment metrics
		entityAnalysis.TotalArticles = len(entityAnalysis.Analyses)

		var articlesList []string
		for _, analysis := range entityAnalysis.Analyses {
			articlesList = append(articlesList, analysis.ArticleMetadata.URL)
		}
		
		// Calculate average sentiment using the pre-collected sentiment scores
		var avgSentiment float32
		scores := entitySentiments[entityName]
		if len(scores) > 0 {
			var sum float32
			for _, score := range scores {
				sum += score
			}
			avgSentiment = sum / float32(len(scores))
		}

		// Generate entity summary
		entityReport, err := s.GenerateEntityReport(articles, entityName, req.UserID)
		if err != nil {
			s.logger.Errorf("Failed to generate summary for entity %s: %v", entityName, err)
			continue
		}

		entityReportResponse := EntityReport{
			EntityName:       entityName,
			Summary:         entityReport.Summary,
			ArticleCount:    entityReport.ArticleCount,
			AverageSentiment: float32(math.Round(float64(avgSentiment)*100) / 100),
			SentimentLabel:   helper.GetSentimentLabel(avgSentiment),
			Articles:		articlesList,
		}
		
		entitiesReportsResponse = append(entitiesReportsResponse, entityReportResponse)
	}

	return &GetMyReportsResponse{
		Entities: entitiesReportsResponse,
	}, nil
}

// Helper function to create Analysis from ArticleData
func createAnalysisFromArticle(article server.ArticleData) Analysis {
	var publishedDate, scrapedAt time.Time
	
	// Safely parse dates
	if article.PublishedDate != "" {
		if parsed, err := time.Parse(time.RFC3339, article.PublishedDate); err == nil {
			publishedDate = parsed
		}
	}
	if article.ScrapedAt != "" {
		if parsed, err := time.Parse(time.RFC3339, article.ScrapedAt); err == nil {
			scrapedAt = parsed
		}
	}

	// Create analysis results with safe initialization
	analysisResults := AnalysisResults{
		Entities: []AnalysisEntity{},
		Topics:   []AnalysisTopic{},
	}

	// Safely convert entities
	if article.Entities != nil {
		for entityName, entity := range article.Entities {
			var sentimentScore float32
			
			sentimentScore = entity.SentimentScore
			
			analysisEntity := AnalysisEntity{
				Name:           entityName,
				RelatedTopics:  entity.RelatedTopics,
				SentimentLabel: helper.GetSentimentLabel(sentimentScore),
				SentimentScore: sentimentScore,
			}
			analysisResults.Entities = append(analysisResults.Entities, analysisEntity)
		}
	}

	// Safely convert topics
	if article.Topics != nil {
		for _, topic := range article.Topics {
			analysisTopic := AnalysisTopic{
				Name:           topic.Name,
				SentimentLabel: topic.SentimentLabel,
				SentimentScore: topic.SentimentScore,
			}
			analysisResults.Topics = append(analysisResults.Topics, analysisTopic)
		}
	}

	return Analysis{
		AnalysisResults: analysisResults,
		ArticleMetadata: ArticleMetadata{
			ArticleSummary: article.ArticleSummary,
			ID:            article.ArticleID,
				PublishedDate: func() string {
					if !publishedDate.IsZero() {
						return publishedDate.Format(time.RFC3339)
					}
					return ""
				}(),
				ScrapedAt: func() string {
					if !scrapedAt.IsZero() {
						return scrapedAt.Format(time.RFC3339)
					}
					return ""
				}(),
				Title: article.Title,
				URLID: article.URLID,
				URL:   article.URL,
		},
	}
}

func getSentimentLabel(score float32) string {
	switch {
	case score > 0:
		return "Positive"
	case score < 0:
		return "Negative"
	default:
		return "Neutral"
	}
}

func (s *reportsApi) GenerateEntityReport(articles []server.ArticleData, entityName string, userID int) (*EntityReport, error) {
    // First, get the entity ID
    var entity entities.Entity
    if err := s.db.Where("name = ?", entityName).First(&entity).Error; err != nil {
        return nil, fmt.Errorf("entity not found: %v", err)
    }

    // Get article IDs and convert server.ArticleData to []articles.Article
    var articleIDs []int
    var relevantArticles []string
    for _, article := range articles {
        if _, exists := article.Entities[entityName]; exists {
            articleIDs = append(articleIDs, article.ArticleID)
            relevantArticles = append(relevantArticles, article.URL)
        }
    }

    // Sort article IDs for consistent checking
    sort.Ints(articleIDs)

    // Check if we have a recent entity report with the same articles
    var existingReport entity_reportsvc.EntityReport
    err := s.db.Preload("Articles").
        Where("entity_reports.entity_id = ?", entity.ID).
        Where("last_analyzed > ?", time.Now().Add(-24*time.Hour)).
        First(&existingReport).Error

    // If we found a recent report (less than 24 hours old)
    if err == nil {
        // Associate report with current user if not already associated
        s.associateReportWithUser(existingReport.ID, userID)

        return &EntityReport{
            EntityName:    entityName,
            Summary:       existingReport.Summary,
            ArticleCount:  existingReport.ArticleCount,
            Articles:      relevantArticles,
        }, nil
    }

    // If we're here, we need to generate a new report
    var summaries []string
    for _, article := range articles {
        if _, exists := article.Entities[entityName]; exists {
            if article.ArticleSummary != "" {
                summaries = append(summaries, article.ArticleSummary)
            }
        }
    }

    if len(summaries) == 0 {
        return nil, fmt.Errorf("no summaries found for entity %s", entityName)
    }

    // Generate new summary using OpenAI
    summary, err := s.generateOpenAISummary(summaries, entityName)
    if err != nil {
        return nil, err
    }

    // Create or update entity report
    newReport := entity_reportsvc.EntityReport{
        EntityID:     entity.ID,
        Summary:      summary,
        ArticleCount: len(summaries),
        LastAnalyzed: time.Now(),
    }

    // Start a transaction
    tx := s.db.Begin()
    if err := tx.Create(&newReport).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create entity report: %v", err)
    }

    // Associate articles with the report
    for _, articleID := range articleIDs {
        if err := tx.Create(&entity_reportsvc.EntityReportArticle{
            EntityReportID: newReport.ID,
            ArticleID:      articleID,
        }).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to associate article: %v", err)
        }
    }

    // Associate with current user
    if err := tx.Create(&entity_reportsvc.UserEntityReport{
        EntityReportID: newReport.ID,
        UserID:         userID,
    }).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to associate user: %v", err)
    }

    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &EntityReport{
        EntityName:    entityName,
        Summary:       summary,
        ArticleCount:  len(summaries),
        Articles:      relevantArticles,
    }, nil
}

// Helper function to generate summary using OpenAI
func (s *reportsApi) generateOpenAISummary(summaries []string, entityName string) (string, error) {
    prompt := fmt.Sprintf(`Based on these %d article summaries about %s, create a comprehensive report:

    Article Summaries:
    %s

    Please provide a concise summary that covers:
    1. Key events and developments
    2. Overall sentiment and public perception
    3. Main relationships and interactions
    4. Notable trends or patterns

    Format the response as a clear, professional summary.`, 
    len(summaries), entityName, strings.Join(summaries, "\n\n"))

    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: "gpt-4o-mini",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    "user",
                    Content: prompt,
                },
            },
            MaxTokens:   1000,
            Temperature: 0.2,
        },
    )

    if err != nil {
        return "", fmt.Errorf("failed to generate report: %w", err)
    }

    if len(resp.Choices) == 0 {
        return "", fmt.Errorf("no response generated from OpenAI")
    }

    return resp.Choices[0].Message.Content, nil
}

// Helper function to associate report with user
func (s *reportsApi) associateReportWithUser(reportID uint, userID int) error {
    // Check if association already exists
    var existing entity_reportsvc.UserEntityReport
    err := s.db.Where("entity_report_id = ? AND user_id = ?", reportID, userID).
        First(&existing).Error

    if err == gorm.ErrRecordNotFound {
        // Create new association
        return s.db.Create(&entity_reportsvc.UserEntityReport{
            EntityReportID: reportID,
            UserID:         userID,
        }).Error
    }

    return err
}
