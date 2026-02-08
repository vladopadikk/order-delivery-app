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
