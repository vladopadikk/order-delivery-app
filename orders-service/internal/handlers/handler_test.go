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
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/handlers/mocks"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

func TestCreateOrderHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockOrderService, userID int64, orderIn models.OrderInput)

	testCases := []struct {
		name                 string
		setUserID            bool
		userID               int64
		inputBody            string
		inputOrder           models.OrderInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			setUserID: true,
			userID:    42,
			inputBody: `{"items":[{"item_id":1,"quantity":2,"price":500}],"delivery_address":"Moscow, Lenina 1"}`,
			inputOrder: models.OrderInput{
				Items:           []models.OrderItemInput{{ProductID: 1, Quantity: 2, Price: 500}},
				DeliveryAddress: "Moscow, Lenina 1",
			},
			mockBehavior: func(m *mocks.MockOrderService, userID int64, orderIn models.OrderInput) {
				m.EXPECT().
					CreateOrder(mock.Anything, userID, orderIn).
					Return(models.OrderResponse{
						OrderID:    52,
						TotalPrice: 1000,
						Status:     models.StatusPaymentPending,
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"order_id":52,"total_price":1000,"status":"PAYMENT_PENDING"}`,
		},
		{
			name:                 "User unauthorized",
			inputBody:            `{"items":[{"item_id":1,"quantity":2,"price":500}],"delivery_address":"Moscow, Lenina 1"}`,
			mockBehavior:         func(m *mocks.MockOrderService, userID int64, orderIn models.OrderInput) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"user unauthorized"}`,
		},
		{
			name:                 "Invalid JSON",
			setUserID:            true,
			userID:               42,
			inputBody:            `{"items:}`,
			inputOrder:           models.OrderInput{},
			mockBehavior:         func(m *mocks.MockOrderService, userID int64, orderIn models.OrderInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid json"}`,
		},
		{
			name:      "Internal server error",
			setUserID: true,
			userID:    42,
			inputBody: `{"items":[{"item_id":1,"quantity":2,"price":500}],"delivery_address":"Moscow, Lenina 1"}`,
			inputOrder: models.OrderInput{
				Items:           []models.OrderItemInput{{ProductID: 1, Quantity: 2, Price: 500}},
				DeliveryAddress: "Moscow, Lenina 1",
			},
			mockBehavior: func(m *mocks.MockOrderService, userID int64, orderIn models.OrderInput) {
				m.EXPECT().
					CreateOrder(mock.Anything, userID, orderIn).
					Return(models.OrderResponse{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			svc := mocks.NewMockOrderService(t)
			handler := NewHandler(svc)

			tC.mockBehavior(svc, tC.userID, tC.inputOrder)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				if tC.setUserID {
					ctx.Set("userID", tC.userID)
				}
				ctx.Next()
			})

			group := r.Group("")
			RegisterRouter(group, handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/orders",
				bytes.NewBufferString(tC.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tC.expectedStatusCode, w.Code)
			assert.Equal(t, tC.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetOrderListHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockOrderService, userID int64)

	testCases := []struct {
		name                 string
		setUserID            bool
		userID               int64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			setUserID: true,
			userID:    42,
			mockBehavior: func(m *mocks.MockOrderService, userID int64) {
				m.EXPECT().
					GetOrderList(mock.Anything, userID).
					Return(models.OrderListResponse{
						UserID: 42,
						Orders: []models.OrderResponse{
							{
								OrderID:    52,
								TotalPrice: 1000,
								Status:     models.StatusPaymentPending,
							},
						},
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":42,"orders":[{"order_id":52,"total_price":1000,"status":"PAYMENT_PENDING"}]}`,
		},
		{
			name:                 "Unauthorized",
			setUserID:            false,
			mockBehavior:         func(m *mocks.MockOrderService, userID int64) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"user unauthorized"}`,
		},
		{
			name:      "Internal server error",
			setUserID: true,
			userID:    42,
			mockBehavior: func(m *mocks.MockOrderService, userID int64) {
				m.EXPECT().
					GetOrderList(mock.Anything, userID).
					Return(models.OrderListResponse{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			svc := mocks.NewMockOrderService(t)
			handler := NewHandler(svc)

			tC.mockBehavior(svc, tC.userID)

			gin.SetMode(gin.TestMode)
			r := gin.New()

			r.Use(func(ctx *gin.Context) {
				if tC.setUserID {
					ctx.Set("userID", tC.userID)
				}
				ctx.Next()
			})

			group := r.Group("")
			RegisterRouter(group, handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/orders", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tC.expectedStatusCode, w.Code)
			assert.Equal(t, tC.expectedResponseBody, w.Body.String())

		})
	}
}
