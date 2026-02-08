package models

type OrderItemInput struct {
	ProductID int64  `json:"item_id"`
	Quantity  uint32 `json:"quantity"`
}

type OrderInput struct {
	Items           []OrderItemInput `json:"items"`
	DeliveryAddress string           `json:"delivery_address"`
}

type OrderResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
