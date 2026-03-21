package models

import "time"

type Order struct {
	ID              int64
	UserID          int64
	Status          string
	TotalPrice      float64
	DeliveryAddress string
	CreatedAt       time.Time
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  uint32
	Price     float64
}
