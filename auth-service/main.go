package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/handlers"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/repository"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	router := gin.Default()

	repo := repository.NewRepository(db)
	service := service.NewService(repo, cfg)
	handler := handlers.NewHandler(service)

	api := router.Group("/api")
	handlers.RegisterRoutes(api, handler)

	router.Run(":" + cfg.AppPort)

}
