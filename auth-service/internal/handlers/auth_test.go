package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/handlers/mocks"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/service"
)

func TestRegisterHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockAuthService, userIn models.UserInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.UserInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "petr", "email": "petr@mail.com", "password": "pass"}`,
			inputUser: models.UserInput{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, userIn models.UserInput) {
				m.EXPECT().
					Register(mock.Anything, userIn).
					Return(models.UserResponse{
						ID:       42,
						Username: "petr",
						Email:    "petr@mail.com",
					}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":42,"username":"petr","email":"petr@mail.com"}`,
		},
		{
			name:                 "Invalid JSON",
			inputBody:            `{"email":}`,
			inputUser:            models.UserInput{},
			mockBehavior:         func(m *mocks.MockAuthService, userIn models.UserInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid json"}`,
		},
		{
			name:      "Email exists",
			inputBody: `{"username": "petr", "email": "petr@mail.com", "password": "pass"}`,
			inputUser: models.UserInput{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, userIn models.UserInput) {
				m.EXPECT().
					Register(mock.Anything, userIn).
					Return(models.UserResponse{}, service.ErrEmailExists)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"email is already registered"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"username": "petr", "email": "petr@mail.com", "password": "pass"}`,
			inputUser: models.UserInput{
				Username: "petr",
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, userIn models.UserInput) {
				m.EXPECT().
					Register(mock.Anything, userIn).
					Return(models.UserResponse{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := mocks.NewMockAuthService(t)
			handler := NewHandler(svc)

			test.mockBehavior(svc, test.inputUser)

			gin.SetMode(gin.TestMode)
			r := gin.New()

			api := r.Group("/api")
			RegisterRoutes(api, handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/register",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestLoginHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockAuthService, loginIn models.LoginInput)

	tests := []struct {
		name                 string
		inputBody            string
		loginInput           models.LoginInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "petr@mail.com", "password": "pass"}`,
			loginInput: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, loginIn models.LoginInput) {
				m.EXPECT().
					Login(mock.Anything, loginIn).
					Return(models.TokenResponse{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token":"access-token","refresh_token":"refresh-token"}`,
		},
		{
			name:      "User not found",
			inputBody: `{"email": "petr@mail.com", "password": "pass"}`,
			loginInput: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, loginIn models.LoginInput) {
				m.EXPECT().
					Login(mock.Anything, loginIn).
					Return(models.TokenResponse{}, service.ErrUserNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"user not found"}`,
		},
		{
			name:      "Invalid password",
			inputBody: `{"email": "petr@mail.com", "password": "pass"}`,
			loginInput: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, loginIn models.LoginInput) {
				m.EXPECT().
					Login(mock.Anything, loginIn).
					Return(models.TokenResponse{}, service.ErrInvalidPassword)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid password"}`,
		},
		{
			name:                 "Invalid JSON",
			inputBody:            `{"email": }`,
			loginInput:           models.LoginInput{},
			mockBehavior:         func(m *mocks.MockAuthService, loginIn models.LoginInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid json"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"email": "petr@mail.com", "password": "pass"}`,
			loginInput: models.LoginInput{
				Email:    "petr@mail.com",
				Password: "pass",
			},
			mockBehavior: func(m *mocks.MockAuthService, loginIn models.LoginInput) {
				m.EXPECT().
					Login(mock.Anything, loginIn).
					Return(models.TokenResponse{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := mocks.NewMockAuthService(t)
			handler := NewHandler(svc)

			test.mockBehavior(svc, test.loginInput)

			gin.SetMode(gin.TestMode)
			r := gin.New()

			api := r.Group("/api")
			RegisterRoutes(api, handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/login",
				bytes.NewBufferString(test.inputBody),
			)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
