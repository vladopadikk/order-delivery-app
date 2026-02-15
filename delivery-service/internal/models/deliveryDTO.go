package models

const (
	StatusPending    = "DELIVERY_PENDING"
	StatusInProgress = "DELIVERY_IN_PROGRESS"
	StatusDelivered  = "DELIVERED"
)

type DeliveryResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
