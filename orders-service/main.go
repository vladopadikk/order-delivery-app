package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/database"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/handlers"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/kafka/consumer"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/kafka/producer"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/middleware"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/outbox"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/repository"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	broker := cfg.KafkaHost + ":" + cfg.KafkaPort
	kafkaProducer := producer.NewProducer(broker)
	defer kafkaProducer.Close()

	repo := repository.NewRepository(db)
	orderService := service.NewService(repo)
	handler := handlers.NewHandler(orderService)

	relay := outbox.NewRelay(repo, kafkaProducer, cfg.OutboxPollInterval, cfg.OutboxBatchSize)
	relayDone := make(chan struct{})
	go func() {
		defer close(relayDone)
		relay.Run(ctx)
	}()

	orderConsumer := consumer.NewConsumer(broker, orderService)
	go orderConsumer.Run(ctx)

	router := gin.Default()
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	handlers.RegisterRouter(protected, handler)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("http shutdown error:", err)
	}

	<-relayDone
}
