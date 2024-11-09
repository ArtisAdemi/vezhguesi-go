package reports

import (
	"time"
)

type CreateReportRequest struct {
	UserID    int       `json:"-"`
	Subject   string    `json:"subject"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type ReportResponse struct {
	Report Report `json:"report"`
	UserID int    `json:"userId"`
	Articles []Articles `json:"articles"`
}

type ReportsResponse struct {
	Reports []Report `json:"reports"`
}

type GetReportsRequest struct {
	UserID int `json:"-"`
	Terms []string `json:"terms"`
}

type IDRequest struct {
	ID     int `json:"-"`
	UserID int `json:"-"`
}

type UpdateReportRequest struct {
	ID         int           `json:"-"`
	UserID     int           `json:"-"`
	Title      string        `json:"title"`
	Subject    string        `json:"subject"`
	ReportText string        `json:"reportText"`
	Entities   []ReportEntity `json:"entities"`
	SourceID   int           `json:"sourceId"`
	Findings   string        `json:"findings"`
	Sentiment  int           `json:"sentiment"`
	StartDate  time.Time     `json:"startDate"`
	EndDate    time.Time     `json:"endDate"`
}

type ReportEntity struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AnalysisEntity struct {
	Name           string   `json:"name"`
	RelatedTopics  []string `json:"related_topics"`
	SentimentLabel string   `json:"sentiment_label"`
	SentimentScore float32  `json:"sentiment_score"`
}

type AnalysisTopic struct {
	Name            string   `json:"name"`
	RelatedEntities []string `json:"related_entities"`
	SentimentLabel  string   `json:"sentiment_label"`
	SentimentScore  float32  `json:"sentiment_score"`
}

type AnalysisResults struct {
	Entities []AnalysisEntity `json:"entities"`
	Topics   []AnalysisTopic  `json:"topics"`
}

type ArticleMetadata struct {
	ArticleSummary string `json:"article_summary"`
	ID             int    `json:"id"`
	PublishedDate  string `json:"published_date"`
	ScrapedAt      string `json:"scraped_at"`
	Title          string `json:"title"`
	URLID          int    `json:"url_id"`
	URL            string `json:"url"`
}

type Analysis struct {
	AnalysisResults AnalysisResults `json:"analysis_results"`
	ArticleMetadata ArticleMetadata `json:"article_metadata"`
	
}

type GetReportsResponse struct {
	Analyses       []Analysis `json:"analyses"`
	TotalArticles  int        `json:"total_articles"`
}

type EntityReport struct {
    EntityName string `json:"entity_name"`
    Summary string `json:"summary"`
    ArticleCount int `json:"article_count"`
	AverageSentiment float32 `json:"average_sentiment"`
    SentimentLabel string `json:"sentiment_label"`
    TimeRange string `json:"time_range"`
	Articles []string `json:"articles"`
}


type Articles struct {
	ID            int `json:"id"`
	ConfigID      int `json:"config_id"`
	URLID         int `json:"url_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	PublishedDate time.Time `json:"published_date"`
	ScrapedAt     time.Time `json:"scraped_at"`
}

// Define the structure for the JSON response
type GetAnalyzesResponse struct {
	Query   Query   `json:"query"`
	Results Results `json:"results"`
}

type Query struct {
	SearchTerms []string `json:"search_terms"`
}

type Results struct {
	Articles     []Article `json:"articles"`
	TotalArticles int       `json:"total_articles"`
}

type Article struct {
	ArticleID       int               `json:"article_id"`
	ArticleSummary  string            `json:"article_summary"`
	Entities        map[string]Entity `json:"entities"`
	EntitySentimentScores []float32   `json:"entity_sentiment_scores"`
	PublishedDate   time.Time         `json:"published_date"`
	ScrapedAt       time.Time         `json:"scraped_at"`
	Title           string            `json:"title"`
	TopicSentimentScores []float32    `json:"topic_sentiment_scores"`
	Topics          map[string]Topic  `json:"topics"`
	URL             string            `json:"url"`
}

type Entity struct {
	Name            string    `json:"name"`
	RelatedTopics   []string  `json:"related_topics"`
	SentimentLabel  string    `json:"sentiment_label"`
	SentimentScore  float32   `json:"sentiment_score"`
}

type Topic struct {
	Name            string   `json:"name"`
	RelatedEntities []string `json:"related_entities"`
	SentimentLabel  string   `json:"sentiment_label"`
	SentimentScore  float32  `json:"sentiment_score"`
}

// New models for GetMyReports response
type EntityAnalysis struct {
    EntityName        string     `json:"entity_name"`
    Analyses         []Analysis  `json:"analyses"`
    TotalArticles    int        `json:"total_articles"`
    AverageSentiment float32    `json:"average_sentiment"`
    SentimentLabel   string     `json:"sentiment_label"`
}

type GetMyReportsResponse struct {
    Entities []EntityReport `json:"entities"`
}