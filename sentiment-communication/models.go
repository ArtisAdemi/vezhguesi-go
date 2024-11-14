package articles

type Articles struct {
	ID            int    `json:"id"`
	ConfigID      int    `json:"config_id"`
	URLID         int    `json:"url_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	URL           string `json:"url"`
	PublishedDate string `json:"published_date"`
	ScrapedAt     string `json:"scraped_at"`
}

//  return jsonify({
// 	'status': 'completed',
// 	'summary': {
// 		'total_requested': len(article_ids),
// 		'retrieved_from_cache': len(existing_article_ids),
// 		'newly_analyzed': len(articles_to_analyze),
// 		'successful': len(results),
// 		'failed': len(errors)
// 	},
// 	'results': results,
// 	'errors': errors

type AnalyzeArticlesResponse struct {
	Status  string        `json:"status"`
	Summary Summary       `json:"summary"`
	Errors  []string      `json:"errors"`
	Results []ArticleData `json:"results"`
}

type Summary struct {
	TotalRequested     int `json:"total_requested"`
	RetrievedFromCache int `json:"retrieved_from_cache"`
	NewlyAnalyzed      int `json:"newly_analyzed"`
	Successful         int `json:"successful"`
	Failed             int `json:"failed"`
}

type GetAnalyzesResponse struct {
	Query   Query   `json:"query"`
	Results Results `json:"results"`
}

type Query struct {
	SearchTerms []string `json:"search_terms"`
}

type Results struct {
	Articles      []ArticleData `json:"articles"`
	TotalArticles int           `json:"total_articles"`
}

type ArticleData struct {
	ArticleID            int               `json:"article_id"`
	ArticleSummary       string            `json:"article_summary"`
	Entities             map[string]Entity `json:"entities"`
	PublishedDate        string            `json:"published_date"`
	ScrapedAt            string            `json:"scraped_at"`
	Title                string            `json:"title"`
	TopicSentimentScores []interface{}     `json:"topic_sentiment_scores"`
	Topics               map[string]Topic  `json:"topics"`
	URL                  string            `json:"url"`
	URLID                int               `json:"url_id"`
}

type Entity struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	RelatedTopics  []string `json:"related_topics"`
	SentimentLabel string   `json:"sentiment_label"`
	SentimentScore float32  `json:"sentiment_score"`
}

type URL struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

type Topic struct {
	Name            string   `json:"name"`
	RelatedEntities []string `json:"related_entities"`
	SentimentLabel  string   `json:"sentiment_label"`
	SentimentScore  float32  `json:"sentiment_score"`
}
