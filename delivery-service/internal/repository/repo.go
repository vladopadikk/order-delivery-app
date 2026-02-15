package repository

import (
	"context"
	"database/sql"

	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, orderID, userID int64, address string, status string) (models.Delivery, error) {
	query := `
		INSERT INTO deliveries (order_id, user_id, address, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, user_id, address, status, created_at
	`

	var delivery models.Delivery

	err := r.DB.QueryRowContext(ctx,
		query,
		orderID,
		userID,
		address,
		status,
	).Scan(
		&delivery.ID,
		&delivery.OrderID,
		&delivery.UserID,
		&delivery.Address,
		&delivery.Status,
		&delivery.CreatedAt,
	)

	return delivery, err

}
func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	query := `
		UPDATE deliveries 
		SET status = $1
		WHERE order_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, status, orderID)
	return err

}

func (r *Repository) GetStatus(ctx context.Context, userID, orderID int64) (models.Delivery, error) {
	query := `
		SELECT order_id, status 
		FROM deliveries
		WHERE user_id = $1 AND order_id = $2
	`

	var delivery models.Delivery
	err := r.DB.QueryRowContext(ctx, query, userID, orderID).Scan(&delivery.OrderID, &delivery.Status)
	return delivery, err
}
