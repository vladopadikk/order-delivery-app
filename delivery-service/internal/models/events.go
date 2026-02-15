package models

type PaymentEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Status     string
}

type DeliveryEvent struct {
	OrderID int64
	UserID  int64
	Status  string
}
