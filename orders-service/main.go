package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/handlers"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/middleware"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/repository"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/service"
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

	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	handlers.RegisterRouter(protected, handler)

	router.Run(":" + cfg.AppPort)

}
