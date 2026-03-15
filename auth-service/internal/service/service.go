package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vladopadikk/order-delivery-app/auth-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists     = errors.New("email is already registered")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type Service struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo, cfg}

}

func (s *Service) Register(ctx context.Context, input models.UserInput) (models.UserResponse, error) {
	_, err := s.repo.GetByEmail(ctx, input.Email)

	if err == nil {
		return models.UserResponse{}, ErrEmailExists
	}

	if err != sql.ErrNoRows {
		return models.UserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}
	createdAt := time.Now()

	id, err := s.repo.Create(ctx, input.Username, input.Email, string(hashedPass), createdAt)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return models.UserResponse{
		ID:       id,
		Username: input.Username,
		Email:    input.Email,
	}, nil
}

func (s *Service) Login(ctx context.Context, loginIn models.LoginInput) (models.TokenResponse, error) {
	u, err := s.repo.GetByEmail(ctx, loginIn.Email)
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("db error: %w", err)
	}

	if u == nil {
		return models.TokenResponse{}, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(loginIn.Password))
	if err != nil {
		return models.TokenResponse{}, ErrInvalidPassword
	}

	accessToken, err := GenerateAccessToken(u.ID, s.cfg.JWTSecret)
	if err != nil {
		return models.TokenResponse{}, err
	}
	refreshToken, err := GenerateRefreshToken(u.ID, s.cfg.JWTSecret)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
