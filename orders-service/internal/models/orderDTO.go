package models

type OrderItemInput struct {
	ProductID int64   `json:"item_id"`
	Quantity  uint32  `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderInput struct {
	Items           []OrderItemInput `json:"items"`
	DeliveryAddress string           `json:"delivery_address"`
}

type OrderResponse struct {
	OrderID    int64   `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type OrderListResponse struct {
	UserID int64           `json:"user_id"`
	Orders []OrderResponse `json:"orders"`
}
