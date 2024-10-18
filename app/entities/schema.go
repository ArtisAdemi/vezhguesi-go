package entities

import "time"

type Entity struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
