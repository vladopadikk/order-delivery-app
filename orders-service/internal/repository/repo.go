package repository

import (
	"context"
	"database/sql"

	"github.com/vladopadikk/order-delivery-app/orders-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) BeginTx(ctx context.Context) (database.Tx, error) {
	return r.DB.BeginTx(ctx, nil)
}

func (r *Repository) Create(ctx context.Context, exec database.Executor, userID int64, status string, totalPrice float64, deliveryAddress string) (models.Order, error) {
	query := `
		INSERT INTO orders (user_id, status, total_price, delivery_address) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, status, total_price, delivery_address, created_at
	`

	var order models.Order

	err := exec.QueryRowContext(
		ctx,
		query,
		userID,
		status,
		totalPrice,
		deliveryAddress,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.DeliveryAddress,
		&order.CreatedAt,
	)

	return order, err
}

func (r *Repository) CreateItems(ctx context.Context, exec database.Executor, order_id int64, orderItem models.OrderItemInput) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4) 
	`
	_, err := exec.ExecContext(ctx, query, order_id, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

func (r *Repository) GetOrders(ctx context.Context, userID int64) ([]models.OrderResponse, error) {
	query := `
		SELECT o.id, o.total_price, o.status 
		FROM orders o 
		WHERE o.user_id = $1 
		ORDER BY o.created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderResponse

	for rows.Next() {
		var order models.OrderResponse
		if err := rows.Scan(
			&order.OrderID,
			&order.TotalPrice,
			&order.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	query := `
		UPDATE orders 
		SET status = $1
		WHERE id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, status, orderID)
	return err

}

func (r *Repository) CreateOutboxMessage(ctx context.Context, exec database.Executor, msg models.OutboxMessage) error {
	query := `
		INSERT INTO outbox (event_id, aggregate_id, topic, message_key, payload)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := exec.ExecContext(ctx, query, msg.EventID, msg.AggregateID, msg.Topic, msg.Key, string(msg.Payload))
	return err
}

func (r *Repository) GetUnpublishedOutbox(ctx context.Context, limit int) ([]models.OutboxMessage, error) {
	query := `
		SELECT id, event_id::text, aggregate_id, topic, message_key, payload::text, created_at
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY id
		LIMIT $1
	`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.OutboxMessage

	for rows.Next() {
		var msg models.OutboxMessage
		var payload string

		if err := rows.Scan(
			&msg.ID,
			&msg.EventID,
			&msg.AggregateID,
			&msg.Topic,
			&msg.Key,
			&payload,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}

		msg.Payload = []byte(payload)
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) MarkOutboxPublished(ctx context.Context, id int64) error {
	query := `
		UPDATE outbox
		SET published_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
