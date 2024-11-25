package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	articlesvc "vezhguesi/app/articles"
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
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	subjectList := strings.Split(req.Subject, ",")

	// First try with article_entities join from local database
	var articles []articlesvc.Article
	err = s.db.
		Joins("JOIN article_entities ON articles.id = article_entities.article_id").
		Where("article_entities.entity_name IN ?", subjectList).
		Preload("EntityRelations").
		Find(&articles).Error

	// If no articles found in local DB, try fetching from server
	if err != nil || len(articles) == 0 {
		serverArticles, err := s.sentiment.FetchArticlesByEntity(subjectList)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch articles by entity: %v", err)
		}
		articles = serverArticles

		// Create entity relations for the newly fetched articles
		for _, article := range articles {
			for _, subject := range subjectList {
				relation := articlesvc.ArticleEntity{
					ArticleID:      article.ID,
					EntityName:     subject,
					SentimentScore: 0,
					SentimentLabel: "neutral",
				}
				
				if err := s.db.Save(&relation).Error; err != nil {
					s.logger.Errorf("Failed to save article-entity relation: %v", err)
				}
			}
		}
	}

	report := &Report{
		Subject:    req.Subject,
		UserID:     req.UserID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}

	if err := s.db.Create(&report).Error; err != nil {
		return nil, fmt.Errorf("failed to create report: %v", err)
	}

	// Convert to response format
	var articlesList []Articles
	for _, article := range articles {
		s.logger.Infof("[Create] Article published date: %v , article id: %v", article.PublishedDate, article.ID)
		if !article.PublishedDate.Before(req.StartDate) && !article.PublishedDate.After(req.EndDate) {
			articlesList = append(articlesList, Articles{
				ID:      article.ID,
				Title:   article.Title,
				Content: article.Content,
			})
		}
	}

	resp := &ReportResponse{
		Report:   *report,
		Articles: articlesList,
	}

	return resp, nil
}

func (s *reportsApi) validateCreateRequest(req *CreateReportRequest) error {
	if req.UserID == 0 {
		return fmt.Errorf("user id is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if req.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}
	if req.EndDate.IsZero() {
		return fmt.Errorf("end date is required")
	}
	return nil
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
	// Log the request
	s.logger.Infof("Getting reports for user ID: %d", req.UserID)

	var reports []Report
	result := s.db.Where("user_id = ?", req.UserID).
		Order("id DESC").
		Find(&reports)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching reports: %v", result.Error)
	}

	// Log the found reports
	s.logger.Infof("Found %d reports", len(reports))
	for _, r := range reports {
		s.logger.Infof("Report subject: %s", r.Subject)
	}

	var terms []string
	for _, report := range reports {
		terms = append(terms, report.Subject)
	}

	// Log the terms we're searching for
	s.logger.Infof("Searching for terms: %v", terms)

	response, err := s.sentiment.GetAnalyzes(terms)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analyzed reports: %v", err)
	}

	// Log the analysis response
	s.logger.Infof("Got analysis response with %d articles", len(response.Results.Articles))

	// Create a map of entities from reports for quick lookup
	requestedEntities := make(map[string]bool)
	var requestedEntityNames []string
	for _, report := range reports {
		entities := strings.Split(report.Subject, ",")
		for _, entity := range entities {
			entityName := strings.TrimSpace(entity)
			s.logger.Infof("Processing requested entity: %s", entityName)
			requestedEntities[strings.ToLower(entityName)] = true
			requestedEntityNames = append(requestedEntityNames, entityName)
		}
	}

	// Group analyses by entity and track sentiment scores
	entityMap := make(map[string]*EntityAnalysis)
	entitySentiments := make(map[string][]float32)
	relatedEntitiesMap := make(map[string]map[string]Entity)
	
	articles := response.Results.Articles
	
	// Track the full names of entities
	entityFullNames := make(map[string]string)
	
	for _, article := range articles {
		
		analysis := createAnalysisFromArticle(article)

		for entityName, entity := range article.Entities {
			s.logger.Infof("Checking entity: %s", entityName)
			
			// Check if this entity matches any of the requested entities
			var matchedRequestedEntity string
			isRequestedEntity := false
			
			for _, requestedName := range requestedEntityNames {
				if isSameEntity(requestedName, entityName) {
					isRequestedEntity = true
					matchedRequestedEntity = requestedName
					
					// Initialize related entities map if not exists
					if _, exists := relatedEntitiesMap[matchedRequestedEntity]; !exists {
						relatedEntitiesMap[matchedRequestedEntity] = make(map[string]Entity)
					}
					
					// Collect all other entities as related
					for otherName, otherEntity := range article.Entities {
						if otherName != entityName {
							relatedEntitiesMap[matchedRequestedEntity][otherName] = Entity{
								Name: otherName,
								Type: otherEntity.Type, // Add type if available
							}
						}
					}
					break
				}
			}
			
			if isRequestedEntity {
				// Use the original entityName instead of matchedRequestedEntity as the key
				if _, exists := entityMap[matchedRequestedEntity]; !exists {
					entityMap[matchedRequestedEntity] = &EntityAnalysis{
						EntityName: matchedRequestedEntity, // Use the original matched name
						Analyses:   []Analysis{},
					}
					entitySentiments[matchedRequestedEntity] = []float32{}
				}
				
				entityMap[matchedRequestedEntity].Analyses = append(entityMap[matchedRequestedEntity].Analyses, analysis)
				entitySentiments[matchedRequestedEntity] = append(entitySentiments[matchedRequestedEntity], entity.SentimentScore)
			}
		}
	}

	// Log the results before returning
	s.logger.Infof("Found %d matching entities", len(entityMap))
	for entityName := range entityMap {
		s.logger.Infof("Matched entity: %s", entityName)
	}

	var entitiesReportsResponse []EntityReport

	// Before processing entities, ensure they exist in the database
	for _, fullName := range entityFullNames {
		var existingEntity entities.Entity
		if err := s.db.Where("name = ?", fullName).First(&existingEntity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create the entity if it doesn't exist
				newEntity := entities.CreateEntityRequest{
					Name: fullName,
					Type: "PERSON", // or appropriate type
				}
				_, err = s.entitiesApi.Create(&newEntity)
				if err != nil {
					s.logger.Errorf("Failed to create entity %s: %v", fullName, err)
					continue
				}
			} else {
				s.logger.Errorf("Error checking entity %s: %v", fullName, err)
				continue
			}
		}
	}

	// Now process the entities as before
	for entityKey, entityAnalysis := range entityMap {
		// Calculate sentiment metrics
		entityAnalysis.TotalArticles = len(entityAnalysis.Analyses)

		var articlesList []string
		for _, analysis := range entityAnalysis.Analyses {
			articlesList = append(articlesList, analysis.ArticleMetadata.URL)
		}
		
		// Calculate average sentiment using the pre-collected sentiment scores
		var avgSentiment float32
		scores := entitySentiments[entityKey]
		if len(scores) > 0 {
			var sum float32
			for _, score := range scores {
				sum += score
			}
			avgSentiment = sum / float32(len(scores))
		}

		// Generate entity summary
		entityReport, err := s.GenerateEntityReport(articles, entityKey, req.UserID)
		if err != nil {
			s.logger.Errorf("Failed to generate summary for entity %s: %v", entityKey, err)
			continue
		}

		entityReportResponse := EntityReport{
			EntityName:        entityKey, // Use entityKey directly instead of entityFullNames[entityKey]
			Summary:          entityReport.Summary,
			ArticleCount:     entityReport.ArticleCount,
			AverageSentiment: float32(math.Round(float64(avgSentiment)*100) / 100),
			SentimentLabel:   helper.GetSentimentLabel(avgSentiment),
			Articles:         articlesList,
			RelatedEntities:  mapToSlice(relatedEntitiesMap[entityKey]),
		}
		
		entitiesReportsResponse = append(entitiesReportsResponse, entityReportResponse)
	}

	// Sort entities by ID in descending order
	sort.Slice(entitiesReportsResponse, func(i, j int) bool {
		return entitiesReportsResponse[i].ArticleCount > entitiesReportsResponse[j].ArticleCount
	})

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

func (s *reportsApi) GenerateEntityReport(articles []server.ArticleData, entityName string, userID int) (*EntityReport, error) {
    // Get the full entity name from database
    var entity entities.Entity
    if err := s.db.Where("name ILIKE ?", "%"+entityName+"%").First(&entity).Error; err != nil {
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
            EntityName:    entity.Name,
            Summary:       existingReport.Summary,
            ArticleCount:  existingReport.ArticleCount,
            Articles:      relevantArticles,
        }, nil
    }

    // If we're here, we need to generate a new report
    var summaries []string
    for _, article := range articles {
        if _, exists := article.Entities[entity.Name]; exists {
            if article.ArticleSummary != "" {
                summaries = append(summaries, article.ArticleSummary)
            }
        }
    }

    if len(summaries) == 0 {
        return nil, fmt.Errorf("no summaries found for entity %s", entity.Name)
    }

    // Generate new summary using OpenAI
    summary, err := s.generateOpenAISummary(summaries, entity.Name)
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
        EntityName:    entity.Name,
        Summary:       summary,
        ArticleCount:  len(summaries),
        Articles:      relevantArticles,
    }, nil
}

// Helper function to generate summary using OpenAI
func (s *reportsApi) generateOpenAISummary(summaries []string, entityName string) (string, error) {
    prompt := fmt.Sprintf(`Bazuar në këto %d përmbledhje artikujsh për %s, krijoni një raport të shkurtër dhe të qartë.

    Përmbledhjet e artikujve:
    %s

    Krijoni një raport me pikat e mëposhtme, duke përdorur tekst të thjeshtë dhe pika të shkurtra:

    Filloni me "Raport i Përmbledhur për [entity]" si titull.
    Pastaj përfshini këto seksione në rend:

    - Ngjarjet kryesore dhe zhvillimet
    - Perceptimi i përgjithshëm dhe opinioni publik
    - Marrëdhëniet kryesore dhe ndërveprimet
    - Trendet ose modelet e dukshme

    E rëndësishme: Përdorni pika të shkurtra dhe mos përdorni asnjë formatim të veçantë.`, 
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
            MaxTokens:   500,  // Reduce max tokens to keep it concise
            Temperature: 0.2,
        },
    )

    if err != nil {
        return "", fmt.Errorf("failed to generate report: %w", err)
    }

    if len(resp.Choices) == 0 {
        return "", fmt.Errorf("no response generated from OpenAI")
    }

    // Clean up the response
    summary := resp.Choices[0].Message.Content

    // Replace bullet points with new lines
    summary = strings.ReplaceAll(summary, "- ", "\n- ")

    // Remove any existing newlines and replace them with spaces
    summary = strings.ReplaceAll(summary, "\\n", " ")
    summary = strings.ReplaceAll(summary, "\n", " ")

    // Remove multiple spaces
    summary = strings.Join(strings.Fields(summary), " ")

    return summary, nil
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

// Helper function to check if an article is already in the slice
func containsArticle(articles []articlesvc.Article, article articlesvc.Article) bool {
	for _, a := range articles {
		if a.ID == article.ID {
			return true
		}
	}
	return false
}

// Helper function to remove duplicate related entities
func uniqueRelatedEntities(entities []Entity) []Entity {
	seen := make(map[string]bool)
	unique := []Entity{}
	
	for _, entity := range entities {
		if !seen[strings.ToLower(entity.Name)] {
			seen[strings.ToLower(entity.Name)] = true
			unique = append(unique, entity)
		}
	}
	
	return unique
}

// Helper function to check if two entity names refer to the same entity
func isSameEntity(reportEntity, articleEntity string) bool {
    reportEntity = strings.ToLower(strings.TrimSpace(reportEntity))
    articleEntity = strings.ToLower(strings.TrimSpace(articleEntity))
    return strings.Contains(articleEntity, reportEntity) || strings.Contains(reportEntity, articleEntity)
}

// Add this helper function
func mapToSlice(entityMap map[string]Entity) []Entity {
    entities := make([]Entity, 0, len(entityMap))
    for _, entity := range entityMap {
        entities = append(entities, entity)
    }
    return entities
}
