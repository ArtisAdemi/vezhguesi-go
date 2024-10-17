package users

import (
	"time"
)

const UserTableName = "users"

type User struct {
	ID            int     `gorm:"primaryKey"`
	Email         string  `gorm:"unique"`
	Username      *string `gorm:"unique"`
	Password      string
	FirstName     string
	LastName      string
	Status        string
	AvatarImgKey  string
	Active        bool
	Phone         string
	VerifiedEmail bool
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}