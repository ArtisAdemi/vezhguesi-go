package entity_reports

import (
	"time"
	"vezhguesi/app/articles"
	"vezhguesi/app/entities"
	"vezhguesi/core/users"
)

type EntityReport struct {
	ID          uint              `gorm:"primaryKey"`
	EntityID    uint              `gorm:"not null"`
	Entity      entities.Entity   `gorm:"foreignKey:EntityID"`
	Summary     string            `gorm:"type:text"`
	
	// Many-to-many relationship with Articles
	Articles    []articles.Article `gorm:"many2many:entity_report_articles;"`
	
	// Many-to-many relationship with Users
	Users       []users.User      `gorm:"many2many:user_entity_reports;"`
	
	// Metadata
	ArticleCount int              `gorm:"not null"`
	LastAnalyzed time.Time        `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Junction table struct for EntityReport-Article relationship
type EntityReportArticle struct {
	EntityReportID uint
	ArticleID      int
	CreatedAt      time.Time
}

// Junction table struct for User-EntityReport relationship
type UserEntityReport struct {
	UserID         int
	EntityReportID uint
	CreatedAt      time.Time
}
