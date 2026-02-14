package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	"github.com/vladopadikk/order-delivery-app/payments-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/repository"
)

var ErrPaymentNotFound = errors.New("payment not found")

type Service struct {
	repo     *repository.Repository
	producer *producer.Producer
}

func NewService(repo *repository.Repository, producer *producer.Producer) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
	}
}

func (s *Service) ProcessPayment(ctx context.Context, event models.OrderCreatedEvent) error {
	var succEvent models.PaymentSuccessEvent

	if rand.Intn(100) < 80 {
		payment, err := s.repo.Create(ctx, event.OrderID, event.UserID, event.TotalPrice, models.StatusPaymentSuccess)
		if err != nil {
			failEvent := models.PaymentFailedEvent{
				OrderID: event.OrderID,
				Status:  models.StatusPaymentFailed,
				Reason:  "internal server error: " + err.Error(),
			}

			if err := s.producer.PublishPaymentFailed(ctx, failEvent); err != nil {
				return fmt.Errorf("publish error: %w", err)
			}

			return fmt.Errorf("internal server error: %w", err)
		}

		succEvent = models.PaymentSuccessEvent{
			OrderID: payment.OrderID,
			Status:  models.StatusPaymentSuccess,
		}

		if err := s.producer.PublishPaymentSuccess(ctx, succEvent); err != nil {
			return fmt.Errorf("publish error: %w", err)
		}

	} else {
		failEvent := models.PaymentFailedEvent{
			OrderID: event.OrderID,
			Status:  models.StatusPaymentFailed,
			Reason:  "payment system error",
		}

		if err := s.producer.PublishPaymentFailed(ctx, failEvent); err != nil {
			return fmt.Errorf("publish error: %w", err)
		}
	}

	return nil
}

func (s *Service) GetOrderStatus(ctx context.Context, userID, orderID int64) (models.PaymentResponse, error) {
	payment, err := s.repo.GetStatus(ctx, userID, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PaymentResponse{}, ErrPaymentNotFound
		}
		return models.PaymentResponse{}, fmt.Errorf("db error: %w", err)
	}

	return models.PaymentResponse{
		OrderID: payment.OrderID,
		Status:  payment.Status,
	}, nil
}
