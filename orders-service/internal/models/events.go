package models

type OrderCreatedEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Address    string
	Status     string
}

type PaymentSuccessEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Address    string
	Status     string
}
