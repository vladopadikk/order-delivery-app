package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladopadikk/order-delivery-app/notification-service/internal/config"
	"github.com/vladopadikk/order-delivery-app/notification-service/internal/consumer"
	"github.com/vladopadikk/order-delivery-app/notification-service/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	broker := cfg.KafkaHost + ":" + cfg.KafkaPort

	notifier := &service.LogNotifier{}
	service := service.NewService(notifier)

	consumer := consumer.NewConsumer(broker, service)

	go consumer.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	cancel()
}
