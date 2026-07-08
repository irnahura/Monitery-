package models

import "time"

type HealthCheck struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	MonitorID      uint      `json:"monitor_id" gorm:"index;not null"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMS int64     `json:"response_time_ms"`
	IsUp           bool      `json:"is_up"`
	CheckedAt      time.Time `json:"checked_at" gorm:"index;not null"`
}
