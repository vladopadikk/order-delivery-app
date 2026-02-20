package models

const (
	StatusPaymentPending = "PAYMENT_PENDING"
	StatusPaid           = "PAID"
	StatusFailed         = "FAILED"
)

type OrderCreatedEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Address    string  `json:"address"`
	Status     string  `json:"status"`
}

type PaymentSuccessEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Address    string  `json:"address"`
	Status     string  `json:"status"`
}

type PaymentFailedEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	Reason     string  `json:"reason"`
}
