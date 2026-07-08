package models

import (
	"time"

	"gorm.io/datatypes"
)

type Monitor struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	UserID             uint           `json:"user_id" gorm:"index;not null"`
	User               User           `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Name               string         `json:"name" gorm:"not null"`
	URL                string         `json:"url" gorm:"not null"`
	IntervalSeconds    int            `json:"interval_seconds" gorm:"not null;default:60"`
	TimeoutSeconds     int            `json:"timeout_seconds" gorm:"not null;default:10"`
	Method             string         `json:"method" gorm:"not null;default:GET"`
	Headers            datatypes.JSON `json:"headers"`
	FollowRedirects    bool           `json:"follow_redirects" gorm:"not null;default:true"`
	ValidateSSL        bool           `json:"validate_ssl" gorm:"not null;default:true"`
	RetryCount         int            `json:"retry_count" gorm:"not null;default:0"`
	LatestStatusCode   *int           `json:"latest_status_code"`
	LatestResponseTime *int64         `json:"latest_response_time_ms"`
	LatestIsUp         *bool          `json:"latest_is_up"`
	LastCheckedAt      *time.Time     `json:"last_checked_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}
