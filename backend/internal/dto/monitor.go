package dto

import "time"

type MonitorRequest struct {
	Name            string            `json:"name" binding:"required"`
	URL             string            `json:"url" binding:"required"`
	IntervalSeconds int               `json:"interval_seconds"`
	TimeoutSeconds  int               `json:"timeout_seconds"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	FollowRedirects *bool             `json:"follow_redirects"`
	ValidateSSL     *bool             `json:"validate_ssl"`
	RetryCount      int               `json:"retry_count"`
}

type MonitorSummary struct {
	ID                 uint              `json:"id"`
	UserID             uint              `json:"user_id"`
	Name               string            `json:"name"`
	URL                string            `json:"url"`
	IntervalSeconds    int               `json:"interval_seconds"`
	TimeoutSeconds     int               `json:"timeout_seconds"`
	Method             string            `json:"method"`
	Headers            map[string]string `json:"headers"`
	FollowRedirects    bool              `json:"follow_redirects"`
	ValidateSSL        bool              `json:"validate_ssl"`
	RetryCount         int               `json:"retry_count"`
	LatestStatusCode   *int              `json:"latest_status_code"`
	LatestResponseTime *int64            `json:"latest_response_time_ms"`
	LatestIsUp         *bool             `json:"latest_is_up"`
	LastCheckedAt      *time.Time        `json:"last_checked_at"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type HealthCheckResponse struct {
	ID             uint      `json:"id"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMS int64     `json:"response_time_ms"`
	IsUp           bool      `json:"is_up"`
	CheckedAt      time.Time `json:"checked_at"`
}

type MonitorHistoryResponse struct {
	History []HealthCheckResponse `json:"history"`
	Summary AnalyticsSummary      `json:"summary"`
}

type MonitorLatestResponse struct {
	Monitor MonitorSummary       `json:"monitor"`
	Latest  *HealthCheckResponse `json:"latest,omitempty"`
	Summary AnalyticsSummary     `json:"summary"`
}
