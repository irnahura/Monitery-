package models

import "time"

type APIKey struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	User      User       `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	KeyHash   string     `json:"-" gorm:"uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}
