package service

import (
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/repository"
)

type Service struct {
	repo     *repository.Repository
	producer *producer.Producer
}

func NewService(repo *repository.Repository, producer *producer.Producer) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
	}
}
