package models

const (
	StatusPaymentPending = "PAYMENT_PENDING"
	StatusPaymentSuccess = "SUCCESS"
	StatusPaymentFailed  = "FAILED"
)

type OrderCreatedEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Status     string
}

type PaymentSuccessEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Status     string
}

type PaymentFailedEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Status     string
	Reason     string
}
