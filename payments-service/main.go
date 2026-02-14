package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/handlers"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/kafka/consumer"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/middleware"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/repository"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/service"
)

func main() {

	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	router := gin.Default()

	broker := cfg.KafkaHost + ":" + cfg.KafkaPort
	producer := producer.NewProducer(broker)
	defer producer.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo, producer)
	handler := handlers.NewHandler(service)

	consumer := consumer.NewConsumer(broker, service)
	go consumer.Run(context.Background())

	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	handlers.RegisterRoutes(protected, handler)

	router.Run(":" + cfg.AppPort)
}
