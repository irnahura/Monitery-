package dto

import "time"

type CreateAPIKeyRequest struct {
	Name      string     `json:"name" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type APIKeySummary struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type APIKeyCreateResponse struct {
	APIKey APIKeySummary `json:"api_key"`
	Key    string        `json:"key"`
}
