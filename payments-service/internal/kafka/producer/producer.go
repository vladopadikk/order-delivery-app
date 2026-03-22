package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/vladopadikk/order-delivery-app/payments-service/internal/models"
)

type Producer struct {
	successWriter *kafka.Writer
	failedWriter  *kafka.Writer
}

func NewProducer(broker string) *Producer {
	successWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "payment_success",
	})

	failedWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "payment_failed",
	})

	return &Producer{
		successWriter: successWriter,
		failedWriter:  failedWriter,
	}
}

func (p *Producer) Close() {
	p.successWriter.Close()
	p.failedWriter.Close()
}

func (p *Producer) PublishPaymentSuccess(ctx context.Context, event models.PaymentSuccessEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.successWriter.WriteMessages(ctx,
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

func (p *Producer) PublishPaymentFailed(ctx context.Context, event models.PaymentFailedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.failedWriter.WriteMessages(ctx,
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
