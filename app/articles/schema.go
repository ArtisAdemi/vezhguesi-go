package articles

import (
	"time"
	// Ensure this path is correct

	"gorm.io/gorm"
)

type ArticleEntity struct {
	ArticleID  int     `gorm:"primaryKey;autoIncrement:false"`
	EntityName string  `gorm:"primaryKey;autoIncrement:false"`
	SentimentScore float32
	SentimentLabel string
}

type Article struct {
	ID            int              `gorm:"primaryKey"`
	ConfigID      int
	URLID         int
	URL           URL              `gorm:"foreignKey:URLID"`
	EntityRelations []ArticleEntity `gorm:"foreignKey:ArticleID"`
	Title         string
	Content       string
	PublishedDate time.Time
	ScrapedAt     time.Time
}

type URL struct {
	ID        int            `gorm:"primaryKey"`
	Path      string         `gorm:"not null" json:"path"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
