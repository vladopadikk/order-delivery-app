package models

import "time"

const (
	StatusPaymentPending = "PAYMENT_PENDING"
	StatusPaymentSuccess = "SUCCESS"
	StatusPaymentFailed  = "FAILED"
)

type Payment struct {
	ID        int64
	OrderID   int64
	UserID    int64
	Amount    float64
	Status    string
	CreatedAt time.Time
}
