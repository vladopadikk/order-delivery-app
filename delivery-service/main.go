package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/handlers"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/kafka/consumer"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/middleware"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/repository"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/service"
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
