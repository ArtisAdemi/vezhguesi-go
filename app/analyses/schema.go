package analyses

import (
	"time"
	"vezhguesi/app/articles"
)

type Analysis struct {
    ID              uint           `gorm:"primaryKey"`
    ArticleID       int            `gorm:"uniqueIndex:idx_article"`
    Article         articles.Article `gorm:"foreignKey:ArticleID"`
    ArticleSummary  string         `gorm:"type:text"`
    Entities        string         `gorm:"type:json"` // Store entities as JSON
    Topics          string         `gorm:"type:json"` // Store topics as JSON
    CreatedAt       time.Time
    UpdatedAt       time.Time
}