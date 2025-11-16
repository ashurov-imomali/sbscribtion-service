package models

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ServiceName string    `gorm:"type:text;not null" json:"service_name"`
	Price       int       `gorm:"not null;check:price>=0" json:"price"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   string    `gorm:"type:varchar(10);not null" json:"start_date"`
	EndDate     *string   `gorm:"type:varchar(10);check:start_date::date < end_date::date" json:"end_date,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type SubscriptionFilter struct {
	UserID      *string
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
}

type Total struct {
	TotalCost *int64 `json:"total_cost"`
}
