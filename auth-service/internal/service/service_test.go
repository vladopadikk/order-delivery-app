package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/service/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterService(t *testing.T) {
	type args struct {
		Username string
		Email    string
		Password string
	}

	type mockBehavior func(m *mocks.MockUserRepository, args args)

	testCases := []struct {
		name             string
		args             args
		mockBehavior     mockBehavior
		expectedResponse models.UserResponse
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "OK",
			args: args{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, args args) {
				m.EXPECT().
					GetByEmail(mock.Anything, args.Email).
					Return(&models.User{}, sql.ErrNoRows)

				m.EXPECT().
					Create(
						mock.Anything,
						args.Username,
						args.Email,
						mock.MatchedBy(func(hash string) bool {
							return bcrypt.CompareHashAndPassword([]byte(hash), []byte("pass")) == nil
						}),
						mock.Anything,
					).Return(42, nil)
			},
			expectedResponse: models.UserResponse{
				ID:       42,
				Username: "petr",
				Email:    "petr@mail.com",
			},
		},
		{
			name: "Email exists",
			args: args{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, args args) {
				m.EXPECT().
					GetByEmail(mock.Anything, args.Email).
					Return(&models.User{
						ID:           42,
						Username:     "petr",
						Email:        "petr@mail.com",
						PasswordHash: "passHash",
					}, nil)
			},
			wantErr:     true,
			expectedErr: ErrEmailExists,
		},
		{
			name: "Create Error",
			args: args{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, args args) {
				m.EXPECT().
					GetByEmail(mock.Anything, args.Email).
					Return(nil, sql.ErrNoRows)

				m.EXPECT().
					Create(
						mock.Anything,
						args.Username,
						args.Email,
						mock.MatchedBy(func(hash string) bool {
							return bcrypt.CompareHashAndPassword([]byte(hash), []byte("pass")) == nil
						}),
						mock.Anything,
					).Return(0, errors.New("connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)
			svc := NewService(repo, &config.Config{JWTSecret: "test-secret"})

			tC.mockBehavior(repo, tC.args)

			got, err := svc.Register(
				context.Background(),
				models.UserInput{
					Username: tC.args.Username,
					Email:    tC.args.Email,
					Password: tC.args.Password,
				},
			)

			if tC.wantErr {
				assert.Error(t, err)
				if tC.expectedErr != nil {
					assert.ErrorIs(t, err, tC.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.expectedResponse, got)
			}

		})
	}
}

func TestLoginService(t *testing.T) {
	type mockBehavior func(m *mocks.MockUserRepository, loginIn models.LoginInput)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)

	testCases := []struct {
		name         string
		loginIn      models.LoginInput
		mockBehavior mockBehavior
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "OK",
			loginIn: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, loginIn models.LoginInput) {
				m.EXPECT().
					GetByEmail(mock.Anything, loginIn.Email).
					Return(&models.User{
						ID:           42,
						Username:     "petr",
						Email:        "petr@mail.com",
						PasswordHash: string(hash),
					}, nil)
			},
		},
		{
			name: "User not found",
			loginIn: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "invalidPass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, loginIn models.LoginInput) {
				m.EXPECT().
					GetByEmail(mock.Anything, loginIn.Email).
					Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: ErrUserNotFound,
		},
		{
			name: "Invalid password",
			loginIn: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "invalidPass",
			},
			mockBehavior: func(m *mocks.MockUserRepository, loginIn models.LoginInput) {
				m.EXPECT().
					GetByEmail(mock.Anything, loginIn.Email).
					Return(&models.User{
						ID:           42,
						Username:     "petr",
						Email:        "petr@mail.com",
						PasswordHash: string(hash),
					}, nil)
			},
			wantErr:     true,
			expectedErr: ErrInvalidPassword,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)
			svc := NewService(repo, &config.Config{JWTSecret: "test-secret"})

			tC.mockBehavior(repo, tC.loginIn)

			got, err := svc.Login(
				context.Background(),
				tC.loginIn,
			)

			if tC.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tC.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.AccessToken)
				assert.NotEmpty(t, got.RefreshToken)
			}
		})
	}
}
