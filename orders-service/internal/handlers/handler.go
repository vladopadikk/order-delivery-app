package handlers

import (
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/service"
)

type Handler struct {
	service *service.Service
}
