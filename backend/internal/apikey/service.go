package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/models"
	"peekaping/backend/internal/repository"
)

type Service struct {
	store *repository.Store
}

func NewService(store *repository.Store) Service {
	return Service{store: store}
}

func (s Service) List(ctx context.Context, userID uint) ([]dto.APIKeySummary, error) {
	keys, err := s.store.ListAPIKeys(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.APIKeySummary, 0, len(keys))
	for _, key := range keys {
		out = append(out, repository.APIKeyToSummary(key))
	}
	return out, nil
}

func (s Service) Create(ctx context.Context, userID uint, input dto.CreateAPIKeyRequest) (dto.APIKeyCreateResponse, error) {
	key, hash, err := newAPIKey()
	if err != nil {
		return dto.APIKeyCreateResponse{}, err
	}
	record := models.APIKey{UserID: userID, Name: input.Name, KeyHash: hash, ExpiresAt: input.ExpiresAt}
	if err := s.store.CreateAPIKey(ctx, &record); err != nil {
		return dto.APIKeyCreateResponse{}, err
	}
	return dto.APIKeyCreateResponse{APIKey: repository.APIKeyToSummary(record), Key: key}, nil
}

func (s Service) Delete(ctx context.Context, userID, id uint) error {
	deleted, err := s.store.DeleteAPIKey(ctx, id, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("api key not found")
	}
	return nil
}

func newAPIKey() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	key := "pk_" + hex.EncodeToString(buf)
	return key, hashKey(key), nil
}

func hashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func ExpiresIn(days int) *time.Time {
	when := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	return &when
}
