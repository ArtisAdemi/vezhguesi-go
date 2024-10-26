package role

import "time"

const (
	RoleTableName           = "roles"
	PermissionTableName     = "permissions"
	RolePermissionTableName = "role_permissions"
)

type Role struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Description *string
	Permissions []Permission `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type Permission struct {
	ID int `gorm:"primaryKey"`
	Name        string
	HTTPMethods []string
	URLs        []string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
