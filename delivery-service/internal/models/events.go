package models

const (
	StatusInProgress = "DELIVERY_IN_PROGRESS"
	StatusDelivered  = "DELIVERED"
)

type PaymentEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Address    string
	Status     string
}

type DeliveryEvent struct {
	OrderID int64
	Status  string
}
