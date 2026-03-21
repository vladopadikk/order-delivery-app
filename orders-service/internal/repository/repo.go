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

	return orders, err
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
