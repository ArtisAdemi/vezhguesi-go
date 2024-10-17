package reports

import "time"

type Report struct {
	ID          string    `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Subject     string    `gorm:"not null"`
	ReportText  string   
	Entities    string
	SourceID    int
	Findings    string
	Sentiment   int 
	StartDate   time.Time 
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
