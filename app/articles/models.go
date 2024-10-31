package articles

import "time"

type Articles struct {
	ID            int `json:"id"`
	ConfigID      int `json:"config_id"`
	URLID         int `json:"url_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	PublishedDate time.Time `json:"published_date"`
	ScrapedAt     time.Time `json:"scraped_at"`
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
	Status string `json:"status"`
	Summary Summary `json:"summary"`
	Errors []string `json:"errors"`
}

type Summary struct {
	TotalRequested     int `json:"total_requested"`
	RetrievedFromCache int `json:"retrieved_from_cache"`
	NewlyAnalyzed      int `json:"newly_analyzed"`
	Successful         int `json:"successful"`
	Failed             int `json:"failed"`
}
