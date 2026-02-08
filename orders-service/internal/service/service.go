package service

import (
	"context"
	"fmt"

	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, orderIn models.OrderInput) (models.OrderResponse, error) {
	tx, err := s.repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	totalPrice := CalculateTotalPrice(orderIn.Items)

	order, err := s.repo.Create(
		ctx,
		tx,
		userID,
		models.StatusPaymentPending,
		totalPrice,
		orderIn.DeliveryAddress,
	)

	if err != nil {
		return models.OrderResponse{}, err
	}

	for _, orderItem := range orderIn.Items {
		err := s.repo.CreateItems(ctx, tx, order.ID, orderItem)
		if err != nil {
			return models.OrderResponse{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.OrderResponse{}, fmt.Errorf("commit tx: %w", err)
	}

	//Тут Kafka

	return models.OrderResponse{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}, err
}

func CalculateTotalPrice(orderItems []models.OrderItemInput) float64 {
	var totalPrice float64

	for _, i := range orderItems {
		totalPrice += i.Price * float64(i.Quantity)
	}

	return totalPrice
}
