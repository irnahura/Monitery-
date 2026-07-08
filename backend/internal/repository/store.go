package repository

import (
	"context"
	"encoding/json"
	"time"

	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *Store) FindUserByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Store) ListMonitorsByUser(ctx context.Context, userID uint) ([]models.Monitor, error) {
	var monitors []models.Monitor
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&monitors).Error; err != nil {
		return nil, err
	}
	return monitors, nil
}

func (s *Store) FindMonitorByIDAndUser(ctx context.Context, id, userID uint) (models.Monitor, error) {
	var monitor models.Monitor
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&monitor).Error; err != nil {
		return models.Monitor{}, err
	}
	return monitor, nil
}

func (s *Store) FindMonitorByID(ctx context.Context, id uint) (models.Monitor, error) {
	var monitor models.Monitor
	if err := s.db.WithContext(ctx).First(&monitor, id).Error; err != nil {
		return models.Monitor{}, err
	}
	return monitor, nil
}

func (s *Store) CreateMonitor(ctx context.Context, monitor *models.Monitor) error {
	return s.db.WithContext(ctx).Create(monitor).Error
}

func (s *Store) UpdateMonitor(ctx context.Context, monitor *models.Monitor, fields map[string]any) error {
	return s.db.WithContext(ctx).Model(monitor).Updates(fields).Error
}

func (s *Store) DeleteMonitorWithChecks(ctx context.Context, monitorID uint, userID uint) (bool, error) {
	deleted := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("monitor_id = ?", monitorID).Delete(&models.HealthCheck{}).Error; err != nil {
			return err
		}
		result := tx.Where("id = ? AND user_id = ?", monitorID, userID).Delete(&models.Monitor{})
		deleted = result.RowsAffected > 0
		return result.Error
	})
	return deleted, err
}

func (s *Store) ListMonitorsForScheduler(ctx context.Context) ([]models.Monitor, error) {
	var monitors []models.Monitor
	if err := s.db.WithContext(ctx).Preload("User").Find(&monitors).Error; err != nil {
		return nil, err
	}
	return monitors, nil
}

func (s *Store) CreateHealthCheck(ctx context.Context, check *models.HealthCheck) error {
	return s.db.WithContext(ctx).Create(check).Error
}

func (s *Store) ListHealthChecks(ctx context.Context, monitorID uint, limit int) ([]models.HealthCheck, error) {
	var checks []models.HealthCheck
	if err := s.db.WithContext(ctx).Where("monitor_id = ?", monitorID).Order("checked_at DESC").Limit(limit).Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Store) LatestHealthCheck(ctx context.Context, monitorID uint) (models.HealthCheck, error) {
	var check models.HealthCheck
	if err := s.db.WithContext(ctx).Where("monitor_id = ?", monitorID).Order("checked_at DESC").First(&check).Error; err != nil {
		return models.HealthCheck{}, err
	}
	return check, nil
}

func (s *Store) CountHealthChecks(ctx context.Context, monitorID uint) (int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.HealthCheck{}).Where("monitor_id = ?", monitorID).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CountUpHealthChecks(ctx context.Context, monitorID uint) (int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.HealthCheck{}).Where("monitor_id = ? AND is_up = ?", monitorID, true).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, userID uint) ([]models.APIKey, error) {
	var keys []models.APIKey
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *Store) CreateAPIKey(ctx context.Context, key *models.APIKey) error {
	return s.db.WithContext(ctx).Create(key).Error
}

func (s *Store) DeleteAPIKey(ctx context.Context, id, userID uint) (bool, error) {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.APIKey{})
	return result.RowsAffected > 0, result.Error
}

func (s *Store) FindAPIKeyByHash(ctx context.Context, keyHash string) (models.APIKey, error) {
	var apiKey models.APIKey
	if err := s.db.WithContext(ctx).Where("key_hash = ?", keyHash).First(&apiKey).Error; err != nil {
		return models.APIKey{}, err
	}
	return apiKey, nil
}

func (s *Store) ValidateAPIKey(ctx context.Context, key string) (uint, bool) {
	apiKey, err := s.FindAPIKeyByHash(ctx, hashKey(key))
	if err != nil {
		return 0, false
	}
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return 0, false
	}
	return apiKey.UserID, true
}

func MonitorToSummary(m models.Monitor) dto.MonitorSummary {
	return dto.MonitorSummary{
		ID:                 m.ID,
		UserID:             m.UserID,
		Name:               m.Name,
		URL:                m.URL,
		IntervalSeconds:    m.IntervalSeconds,
		TimeoutSeconds:     m.TimeoutSeconds,
		Method:             m.Method,
		Headers:            decodeHeaders(m.Headers),
		FollowRedirects:    m.FollowRedirects,
		ValidateSSL:        m.ValidateSSL,
		RetryCount:         m.RetryCount,
		LatestStatusCode:   m.LatestStatusCode,
		LatestResponseTime: m.LatestResponseTime,
		LatestIsUp:         m.LatestIsUp,
		LastCheckedAt:      m.LastCheckedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func HealthCheckToResponse(check models.HealthCheck) dto.HealthCheckResponse {
	return dto.HealthCheckResponse{
		ID:             check.ID,
		StatusCode:     check.StatusCode,
		ResponseTimeMS: check.ResponseTimeMS,
		IsUp:           check.IsUp,
		CheckedAt:      check.CheckedAt,
	}
}

func APIKeyToSummary(key models.APIKey) dto.APIKeySummary {
	return dto.APIKeySummary{
		ID:        key.ID,
		UserID:    key.UserID,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
		ExpiresAt: key.ExpiresAt,
	}
}

func decodeHeaders(raw datatypes.JSON) map[string]string {
	headers := map[string]string{}
	if len(raw) == 0 {
		return headers
	}
	_ = json.Unmarshal(raw, &headers)
	return headers
}
