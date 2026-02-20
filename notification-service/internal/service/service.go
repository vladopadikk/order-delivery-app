package service

import (
	"context"
	"fmt"
	"log"

	"github.com/vladopadikk/order-delivery-app/notification-service/internal/models"
)

type Notifier interface {
	Send(userID int64, message string) error
}

type LogNotifier struct{}

func (l *LogNotifier) Send(userID int64, msg string) error {
	log.Printf("[NOTIFY] user %d: %s", userID, msg)
	return nil
}

type Service struct {
	notifier Notifier
}

func NewService(notifier Notifier) *Service {
	return &Service{
		notifier: notifier,
	}
}

func (s *Service) HandlePaymentSuccess(ctx context.Context, event models.PaymentSuccessEvent) error {
	msg := fmt.Sprintf(
		"Ваш заказ #%d оплачен. Мы передали его в доставку.",
		event.OrderID,
	)

	return s.notifier.Send(event.UserID, msg)
}

func (s *Service) HandlePaymentFailed(ctx context.Context, event models.PaymentFailedEvent) error {
	msg := fmt.Sprintf(
		"Не удалось оплатить заказ #%d. Причина: %s",
		event.OrderID,
		event.Reason,
	)

	return s.notifier.Send(event.UserID, msg)
}

func (s *Service) HandleDelivery(ctx context.Context, event models.DeliveryEvent) error {
	msg := fmt.Sprintf(
		"Ваш заказ #%d доставлен. Спасибо за покупку!",
		event.OrderID,
	)

	return s.notifier.Send(event.UserID, msg)
}
