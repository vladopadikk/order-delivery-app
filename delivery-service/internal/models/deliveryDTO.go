package models

type DeliveryResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
