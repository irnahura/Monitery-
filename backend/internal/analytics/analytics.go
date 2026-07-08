package analytics

import (
	"peekaping/backend/internal/models"

	"gorm.io/gorm"
)

type Summary struct {
	AvailabilityPercent float64 `json:"availability_percent"`
	SLAPercent          float64 `json:"sla_percent"`
}

func Calculate(db *gorm.DB, monitorID uint) (Summary, error) {
	var total int64
	if err := db.Model(&models.HealthCheck{}).Where("monitor_id = ?", monitorID).Count(&total).Error; err != nil {
		return Summary{}, err
	}
	if total == 0 {
		return Summary{AvailabilityPercent: 100, SLAPercent: 100}, nil
	}

	var up int64
	if err := db.Model(&models.HealthCheck{}).Where("monitor_id = ? AND is_up = ?", monitorID, true).Count(&up).Error; err != nil {
		return Summary{}, err
	}
	percent := float64(up) / float64(total) * 100
	return Summary{AvailabilityPercent: percent, SLAPercent: percent}, nil
}
