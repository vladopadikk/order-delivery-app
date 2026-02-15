package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/repository"
)

var ErrDeliveryNotFound = errors.New("delivery not found")

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

func (s *Service) StartDelivery(ctx context.Context, event models.PaymentEvent) error {
	delivery, err := s.repo.Create(
		ctx,
		event.OrderID,
		event.UserID,
		event.Address,
		models.StatusInProgress,
	)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}

	time.Sleep(10 * time.Second)

	if err := s.repo.UpdateStatus(ctx, delivery.OrderID, models.StatusDelivered); err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}

	delEvent := models.DeliveryEvent{
		OrderID: delivery.OrderID,
		Status:  models.StatusDelivered,
	}

	if err = s.producer.PublishDeliveryEvent(ctx, delEvent); err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}

func (s *Service) GetDeliveryStatus(ctx context.Context, userID, orderID int64) (models.DeliveryResponse, error) {
	delivery, err := s.repo.GetStatus(ctx, userID, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DeliveryResponse{}, ErrDeliveryNotFound
		}
		return models.DeliveryResponse{}, fmt.Errorf("db error: %w", err)
	}

	return models.DeliveryResponse{
		OrderID: delivery.OrderID,
		Status:  delivery.Status,
	}, nil
}
