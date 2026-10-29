package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/vladopadikk/order-delivery-app/orders-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) (database.Tx, error)
	Create(ctx context.Context, exec database.Executor, userID int64, status string, totalPrice float64, deliveryAddress string) (models.Order, error)
	CreateItems(ctx context.Context, exec database.Executor, order_id int64, orderItem models.OrderItemInput) error
	CreateOutboxMessage(ctx context.Context, exec database.Executor, msg models.OutboxMessage) error
	GetOrders(ctx context.Context, userID int64) ([]models.OrderResponse, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error
}

type Service struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, orderIn models.OrderInput) (models.OrderResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
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

	// Событие пишется в outbox в той же транзакции, что и заказ.
	// В Kafka его отправит relay отдельно.
	eventID, err := models.NewEventID()
	if err != nil {
		return models.OrderResponse{}, err
	}

	event := models.OrderCreatedEvent{
		EventID:    eventID,
		OccurredAt: time.Now().UTC(),
		OrderID:    order.ID,
		UserID:     userID,
		TotalPrice: order.TotalPrice,
		Address:    order.DeliveryAddress,
		Status:     order.Status,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("marshal order created event: %w", err)
	}

	err = s.repo.CreateOutboxMessage(ctx, tx, models.OutboxMessage{
		EventID:     eventID,
		AggregateID: order.ID,
		Topic:       models.TopicOrderCreated,
		Key:         strconv.FormatInt(order.ID, 10),
		Payload:     payload,
	})
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("create outbox message: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.OrderResponse{}, fmt.Errorf("commit tx: %w", err)
	}

	return models.OrderResponse{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	return s.repo.UpdateStatus(ctx, orderID, status)
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
