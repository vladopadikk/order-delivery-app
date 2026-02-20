package models

const (
	StatusInProgress = "DELIVERY_IN_PROGRESS"
	StatusDelivered  = "DELIVERED"
)

type PaymentEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Address    string  `json:"address"`
	Status     string  `json:"status"`
}

type DeliveryEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Status  string `json:"status"`
}
