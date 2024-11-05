package articles

import (
	"time"
	// Ensure this path is correct
	entitiesvc "vezhguesi/app/entities"

	"gorm.io/gorm"
)

type Article struct {
	ID                     int              `gorm:"primaryKey"`
	ConfigID               int
	URLID                  int
	URL                    URL `gorm:"foreignKey:URLID"`
	Entities               []entitiesvc.Entity `gorm:"many2many:article_entities;"` // Many-to-many relationship
	Title                  string
	Content                string
	PublishedDate          time.Time
	ScrapedAt              time.Time
}

type URL struct {
	ID        int            `gorm:"primaryKey"`
	Path      string         `gorm:"not null" json:"path"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
