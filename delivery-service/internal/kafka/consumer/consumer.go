package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/service"
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

			var event models.PaymentEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json error:", err)
				continue
			}

		}
	}
}
