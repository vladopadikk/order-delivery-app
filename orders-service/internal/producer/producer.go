package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type OrderCreatedEvent struct {
	OrderID    int64
	UserID     int64
	TotalPrice float64
	Status     string
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "order_created",
	})

	return &Producer{writer}
}

func (p *Producer) Close() {
	p.writer.Close()
}

func (p *Producer) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
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
