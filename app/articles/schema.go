package articles

import "time"

type Article struct {
	ID            int       `gorm:"primaryKey"`
	ConfigID      int
	URLID         int
	Title         string
	Content       string
	PublishedDate time.Time
	ScrapedAt     time.Time
}