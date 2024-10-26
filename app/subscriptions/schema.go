package subscriptions

import "time"

const (
	SubscriptionTableName = "subscriptions"
	FeatureTableName      = "features"
)

type Subscription struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Description  string
	Price        float32 `gorm:"precision:10, scale:2"`
	Currency     string
	DurationType string
	DurationTime int
	Features     []Feature `gorm:"foreignKey:SubscriptionID"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

type Feature struct {
	ID             int `gorm:"primaryKey"`
	SubscriptionID int
	Key            string
	Value          string
	CreatedAt      time.Time
}
