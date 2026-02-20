package service

import (
	"context"
	"testing"

	"github.com/vladopadikk/order-delivery-app/notification-service/internal/models"
)

type MockNotifier struct {
	Called  bool
	UserID  int64
	Message string
	Err     error
}

func (m *MockNotifier) Send(userID int64, message string) error {
	m.Called = true
	m.UserID = userID
	m.Message = message
	return m.Err
}

func TestServiceHandlers(t *testing.T) {
	data := []struct {
		name            string
		callHandler     func(s *Service) error
		expectedUserID  int64
		expectedMessage string
		expectError     bool
	}{
		{
			name: "payment success",
			callHandler: func(s *Service) error {
				event := models.PaymentSuccessEvent{
					OrderID: 123,
					UserID:  42,
				}
				return s.HandlePaymentSuccess(context.Background(), event)
			},
			expectedUserID:  42,
			expectedMessage: "Ваш заказ #123 оплачен. Мы передали его в доставку.",
		},
		{
			name: "payment failed",
			callHandler: func(s *Service) error {
				event := models.PaymentFailedEvent{
					OrderID: 123,
					UserID:  42,
					Reason:  "недостаточно средств",
				}
				return s.HandlePaymentFailed(context.Background(), event)
			},
			expectedUserID:  42,
			expectedMessage: "Не удалось оплатить заказ #123. Причина: недостаточно средств",
		},
		{
			name: "delivery completed",
			callHandler: func(s *Service) error {
				event := models.DeliveryEvent{
					OrderID: 123,
					UserID:  42,
				}
				return s.HandleDelivery(context.Background(), event)
			},
			expectedUserID:  42,
			expectedMessage: "Ваш заказ #123 доставлен. Спасибо за покупку!",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			mock := &MockNotifier{}
			service := NewService(mock)

			err := d.callHandler(service)

			if d.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !d.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !mock.Called {
				t.Fatal("notifier was not called")
			}

			if mock.UserID != d.expectedUserID {
				t.Fatalf("expected userID %d, got %d", d.expectedUserID, mock.UserID)
			}

			if mock.Message != d.expectedMessage {
				t.Fatalf("expected message %q, got %q", d.expectedMessage, mock.Message)
			}
		})
	}
}
