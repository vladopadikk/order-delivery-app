package models

type PaymentResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

type PaymentsListResponse struct {
	UserID   int64             `json:"user_id"`
	Payments []PaymentResponse `json:"payments"`
}
