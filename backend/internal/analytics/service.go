package analytics

import (
	"context"

	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/repository"
)

type Service struct {
	store *repository.Store
}

func NewService(store *repository.Store) Service {
	return Service{store: store}
}

func (s Service) Summary(ctx context.Context, monitorID uint) (dto.AnalyticsSummary, error) {
	total, err := s.store.CountHealthChecks(ctx, monitorID)
	if err != nil {
		return dto.AnalyticsSummary{}, err
	}
	if total == 0 {
		return dto.AnalyticsSummary{AvailabilityPercent: 100, SLAPercent: 100}, nil
	}

	up, err := s.store.CountUpHealthChecks(ctx, monitorID)
	if err != nil {
		return dto.AnalyticsSummary{}, err
	}
	percent := float64(up) / float64(total) * 100
	return dto.AnalyticsSummary{AvailabilityPercent: percent, SLAPercent: percent}, nil
}
