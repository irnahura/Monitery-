package auth

import (
	"context"
	"errors"
	"strings"

	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/models"
	"peekaping/backend/internal/repository"
)

type Service struct {
	store     *repository.Store
	jwtSecret string
}

func NewService(store *repository.Store, jwtSecret string) Service {
	return Service{store: store, jwtSecret: jwtSecret}
}

func (s Service) Register(ctx context.Context, input dto.AuthRequest) (dto.AuthResponse, error) {
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	user := models.User{
		Name:         input.Name,
		Email:        strings.ToLower(input.Email),
		PasswordHash: passwordHash,
	}
	if user.Name == "" {
		user.Name = user.Email
	}
	if err := s.store.CreateUser(ctx, &user); err != nil {
		return dto.AuthResponse{}, err
	}
	token, err := IssueJWT(user.ID, s.jwtSecret)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{Token: token, User: toSummary(user)}, nil
}

func (s Service) Login(ctx context.Context, input dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.store.FindUserByEmail(ctx, strings.ToLower(input.Email))
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if !VerifyPassword(user.PasswordHash, input.Password) {
		return dto.AuthResponse{}, errors.New("invalid email or password")
	}
	token, err := IssueJWT(user.ID, s.jwtSecret)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{Token: token, User: toSummary(user)}, nil
}

func (s Service) Profile(ctx context.Context, userID uint) (dto.UserSummary, error) {
	user, err := s.store.FindUserByID(ctx, userID)
	if err != nil {
		return dto.UserSummary{}, err
	}
	return toSummary(user), nil
}

func toSummary(user models.User) dto.UserSummary {
	return dto.UserSummary{ID: user.ID, Name: user.Name, Email: user.Email}
}
