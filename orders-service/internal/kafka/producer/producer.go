package producer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Producer отправляет готовые сообщения в Kafka. Пользуется им только relay.
type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
	})

	return &Producer{writer}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// Publish отправляет сообщение в указанный топик (топик берётся из строки outbox).
func (p *Producer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	err := p.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}
