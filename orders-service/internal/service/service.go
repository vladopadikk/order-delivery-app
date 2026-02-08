package service

import "github.com/vladopadikk/order-delivery-app/orders-service/internal/repository"

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo}
}
