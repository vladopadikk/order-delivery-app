package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/models"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "delivery_completed",
	})

	return &Producer{writer}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) PublishDeliveryEvent(ctx context.Context, event models.DeliveryEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(strconv.FormatInt(event.OrderID, 10)),
			Value: payload,
		},
	)

	if err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return err
}
