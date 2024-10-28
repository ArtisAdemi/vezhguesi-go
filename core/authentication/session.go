package authentication

import "time"

type Session struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	SessionToken string `gorm:"unique;not null"`
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

