package repository

import (
	"context"
	"database/sql"

	"github.com/vladopadikk/order-delivery-app/payments-service/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, orderID, userID int64, amount float64, status string) (models.Payment, error) {
	query := `
		INSERT INTO payments (order_id, user_id, amount, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, user_id, amount, status, created_at
	`

	var payment models.Payment

	err := r.DB.QueryRowContext(ctx,
		query,
		orderID,
		userID,
		amount,
		status,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
	)

	return payment, err

}

func (r *Repository) GetStatus(ctx context.Context, userID, orderID int64) (models.Payment, error) {
	query := `
		SELECT order_id, status 
		FROM payments
		WHERE user_id = $1 AND order_id = $2
	`

	var payment models.Payment
	err := r.DB.QueryRowContext(ctx, query, userID, orderID).Scan(&payment.OrderID, &payment.Status)
	return payment, err
}
