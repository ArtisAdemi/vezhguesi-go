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
}

type ReportsResponse struct {
	Reports []Report `json:"reports"`
}

type GetReportsRequest struct {
	UserID int `json:"-"`
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
