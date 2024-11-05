package articles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	articlesvc "vezhguesi/app/articles"
	entitiesvc "vezhguesi/app/entities"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type serverApi struct {
	db *gorm.DB
	logger log.AllLogger
}

type ServerAPI interface {
	FetchArticles() ([]articlesvc.Article, error)
	AnalyzeArticles(articleIds *[]int) (res *AnalyzeArticlesResponse, err error)
	GetAnalyzes(req []string) (res *GetAnalyzesResponse, err error)
}

func NewServerAPI(db *gorm.DB, logger log.AllLogger) ServerAPI {
	return &serverApi{db: db, logger: logger}
}

func (s *serverApi) FetchArticles() ([]articlesvc.Article, error) {
	// Make an HTTP GET request to fetch the articles data
	resp, err := http.Get(fmt.Sprintf("%s:%s/articles", os.Getenv("SERVER_URL"), os.Getenv("SERVER_ARTICLES_PORT")))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %v", err)
	}
	defer resp.Body.Close()
	var resArticles []articlesvc.Article

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch articles: status code %d", resp.StatusCode)
	}

	// Decode the JSON response
	var articles []Articles
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles data: %v", err)
	}
	for _, article := range articles {
		// Parse the time strings
		scrapedAt, err := time.Parse("2006-01-02T15:04:05.999999", article.ScrapedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse scraped at date: %v", err)
		}

		publishedDate, err := time.Parse("2006-01-02T15:04:05.999999", article.PublishedDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse published date: %v", err)
		}

		// Check if the URL exists
		var url URL
		if err := s.db.Table("urls").Where("path = ?", article.URL).First(&url).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				url.Path = article.URL
				if err := s.db.Table("urls").Create(&url).Error; err != nil {
					return nil, fmt.Errorf("failed to create URL: %v", err)
				}
			} else {
				return nil, fmt.Errorf("failed to query URL: %v", err)
			}
		}

		article.URLID = url.ID

		newArticle := articlesvc.Article{
			ID:            article.ID,
			ConfigID:      article.ConfigID,
			URLID:         article.URLID,
			Title:         article.Title,
			Content:       article.Content,
			PublishedDate: publishedDate,
			ScrapedAt:     scrapedAt,
		}

		// Check if the article already exists
		var existingArticle articlesvc.Article
		if err := s.db.Where("id = ?", article.ID).First(&existingArticle).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&newArticle).Error; err != nil {
					return nil, fmt.Errorf("failed to create article: %v", err)
				}
			} else {
				return nil, fmt.Errorf("failed to query article: %v", err)
			}
		}

		resArticles = append(resArticles, existingArticle)
	}

	return resArticles, nil
}

func (s *serverApi) AnalyzeArticles(articleIds *[]int) (res *AnalyzeArticlesResponse, err error) {
	// Create a map to hold the JSON payload
	payload := map[string]interface{}{
		"article_id": articleIds,
	}

	// Convert the payload to JSON
	articleIdsJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal articleIds: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s:%s/analyze-batch", os.Getenv("SERVER_URL"), os.Getenv("SERVER_ANALYSIS_PORT")), bytes.NewBuffer(articleIdsJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type and authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", os.Getenv("SERVER_API_KEY"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze articles: %v", err)
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to analyze articles: status code %d", resp.StatusCode)
	}

	// Attempt to decode the JSON response
	var response AnalyzeArticlesResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode analyzed articles data: %v", err)
	}

	return &response, nil
}

func (s *serverApi) GetAnalyzes(req []string) (res *GetAnalyzesResponse, err error) {
	baseUrl := fmt.Sprintf("%s:%s/search", os.Getenv("SERVER_URL"), os.Getenv("SERVER_ANALYSIS_PORT"))

	// Create a URL object
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Add query parameters
	query := u.Query()
	for _, item := range req {
		query.Add("terms[]", item) // Use "terms[]" as the parameter name
	}
	u.RawQuery = query.Encode()

	// Create a new HTTP request
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the X-API-Key header
	request.Header.Set("X-API-Key", os.Getenv("SERVER_API_KEY"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get analyzes: status code %d", resp.StatusCode)
	}

	// Decode the JSON response into the GetAnalyzesResponse struct
	var response GetAnalyzesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode analyzed articles data: %v", err)
	}

	// Check and save entities in the database
	for _, articleData := range response.Results.Articles {
		for _, entity := range articleData.Entities {
			var existingEntity entitiesvc.Entity
			if err := s.db.Where("name = ?", entity.Name).First(&existingEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Entity does not exist, create it
					newEntity := entitiesvc.Entity{
						Name:           entity.Name,
						RelatedTopics:  marshalToJson(entity.RelatedTopics),
						SentimentLabel: entity.SentimentLabel,
						SentimentScores: marshalToJson(entity.SentimentScores),
					}
					if err := s.db.Create(&newEntity).Error; err != nil {
						return nil, fmt.Errorf("failed to create entity: %v", err)
					}
				} else {
					return nil, fmt.Errorf("failed to query entity: %v", err)
				}
			}
		}

		url := articleData.URL
		var existingUrl articlesvc.URL
		if err := s.db.Where("path = ?", url).First(&existingUrl).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					newUrl := articlesvc.URL{Path: url}
					if err := s.db.Create(&newUrl).Error; err != nil {
					return nil, fmt.Errorf("failed to create URL: %v", err)
				}
			}
		}
	}

	return &response, nil
}

// Helper function to marshal data to JSON string
func marshalToJson(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "[]"
	}
	return string(jsonData)
}
