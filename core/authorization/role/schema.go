package role

import (
	"time"
)

type Role struct {
	ID          int          `gorm:"primaryKey"`
OrgID *int
	Name        string       `gorm:"not null"`
	Description *string      `gorm:"type:text"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type Permission struct {
	ID          int        `gorm:"primaryKey"`
	Name        string
	HTTPMethods string
	Path        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
