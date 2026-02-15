package models

import "time"

type Delivery struct {
	ID          int64
	OrderID     int64
	UserID      int64
	Address     string
	Status      string
	CreatedAt   time.Time
	CompletedAt time.Time
}
