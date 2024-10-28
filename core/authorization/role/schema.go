package role

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type StringArray []string

// Implement the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Implement the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

type Role struct {
	ID          int          `gorm:"primaryKey"`
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
	HTTPMethods StringArray `gorm:"type:json"` // Use the custom type
	URLs        StringArray `gorm:"type:json"` // Use the custom type
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
