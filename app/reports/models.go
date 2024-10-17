package reports

import "time"

type CreateReportRequest struct {
	UserID    int       `json:"user_id"`
	Subject   string    `json:"subject"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type ReportResponse struct {
	Report Report `json:"report"`
}
