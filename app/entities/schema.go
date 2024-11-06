package entities

import (
	"time"
)

type Entity struct {
	ID             uint              `gorm:"primaryKey"`
	Name           string            `gorm:"not null"`
	Type           string
	RelatedTopics  string            `gorm:"type:json"` // Serialize to JSON
	SentimentLabel string
	SentimentScore float32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

