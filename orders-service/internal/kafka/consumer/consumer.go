package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/service"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.Service
}

func NewConsumer(broker string, service *service.Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "payment_success",
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("consumer stopped")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("kafka error:", err)
				continue
			}

			var event models.PaymentSuccessEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json error:", err)
				continue
			}

			if err := c.service.UpdateOrderStatus(ctx, event); err != nil {
				log.Println("process error:", err)
			}
		}
	}
}
