package orgs

import (
	"time"

	"gorm.io/gorm"
)

const (
	OrgTableName         = "orgs"
	UserOrgRoleTableName = "user_org_roles"
)

type Org struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Slug      string `gorm:"unique;not null"`
	Size      string `gorm:"not null"`
	UserOrgRole []UserOrgRole `gorm:"foreignKey:OrgID"`
	SubscriptionID int
	Subscription Subscription `gorm:"foreignKey:SubscriptionID"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}


type UserOrgRole struct {
	UserID int `gorm:"foreignKey:ID"`
	User User
	OrgID  int
	RoleID int `gorm:"foreignKey:ID"`
	Role Role
	Status string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type Role struct {
	gorm.Model
}

type User struct {
	gorm.Model
}

type Subscription struct {
	gorm.Model
}
