package models

type PaymentEvent struct {
	OrderID int64
	Status  string
}

type DeliveryEvent struct {
	OrderID int64
	Status  string
}
