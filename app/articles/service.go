package articles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type articlesApi struct {
	db *gorm.DB
	logger log.AllLogger
}

type ArticlesAPI interface {
	FetchArticles() ([]Article, error)
	AnalyzeArticles(articleIds *[]int) (res *AnalyzeArticlesResponse, err error)
}

func NewArticlesAPI(db *gorm.DB, logger log.AllLogger) ArticlesAPI {
	return &articlesApi{db: db, logger: logger}
}

func (s *articlesApi) FetchArticles() ([]Article, error) {
	// Make an HTTP GET request to fetch the articles data
	resp, err := http.Get(fmt.Sprintf("%s:%s/articles", os.Getenv("SERVER_URL"), os.Getenv("SERVER_ARTICLES_PORT")))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch articles: status code %d", resp.StatusCode)
	}

	// Decode the JSON response
	var articles []Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles data: %v", err)
	}

	for _, article := range articles {
		_ = s.db.Create(&article)
	}

	return articles, nil
}

func (s *articlesApi) AnalyzeArticles(articleIds *[]int) (res *AnalyzeArticlesResponse, err error) {
	// Create a map to hold the JSON payload
	payload := map[string]interface{}{
		"article_id": articleIds,
	}

	// Convert the payload to JSON
	articleIdsJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal articleIds: %v", err)
	}

	fmt.Println("_____________________________________________________")
	fmt.Println("Sending articleIds to analyze:", string(articleIdsJSON))
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

	fmt.Println("_____________________________________________________")
	fmt.Println("Response from analyze articles:", resp)

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
