package reports

import (
	"time"
	"vezhguesi/app/entities" // Import the entities package
	"vezhguesi/core/users"
)

type Report struct {
	ID          uint              `gorm:"primaryKey"`
	Title       string            `gorm:"not null"`
	Subject     string            `gorm:"not null"`
	UserID      int               
	User        users.User         `gorm:"foreignKey:UserID"`
	ReportText  string   
	Entities    []entities.Entity `gorm:"many2many:report_entities;"` // Updated to use a many-to-many relationship
	SourceID    int
	Findings    string
	Sentiment   int 
	StartDate   time.Time 
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
