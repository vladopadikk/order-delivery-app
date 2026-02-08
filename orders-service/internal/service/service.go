package service

import (
	"context"
	"fmt"

	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/producer"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/repository"
)

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

	event := producer.OrderCreatedEvent{
		OrderID:    order.ID,
		UserID:     userID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}

	if err := s.producer.PublishOrderCreated(ctx, event); err != nil {
		return models.OrderResponse{}, err
	}

	return models.OrderResponse{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}, err
}

func (s *Service) GetOrderList(ctx context.Context, userID int64) (models.OrderListResponse, error) {
	orders, err := s.repo.GetOrders(ctx, userID)
	if err != nil {
		return models.OrderListResponse{}, err
	}

	return models.OrderListResponse{
		UserID: userID,
		Orders: orders,
	}, nil
}

func CalculateTotalPrice(orderItems []models.OrderItemInput) float64 {
	var totalPrice float64

	for _, i := range orderItems {
		totalPrice += i.Price * float64(i.Quantity)
	}

	return totalPrice
}
