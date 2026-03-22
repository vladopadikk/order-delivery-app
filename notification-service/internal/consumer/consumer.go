package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/vladopadikk/order-delivery-app/notification-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/notification-service/internal/service"
)

type Consumer struct {
	paymentSuccReader *kafka.Reader
	paymentFailReader *kafka.Reader
	deliveryReader    *kafka.Reader
	service           *service.Service
}

func NewConsumer(broker string, service *service.Service) *Consumer {
	paymentSuccReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "payment_success",
		GroupID: "notification-service-payment-success",
	})
	paymentFailReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "payment_failed",
		GroupID: "notification-service-payment-failed",
	})
	deliveryReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "delivery_completed",
		GroupID: "notification-service-delivery-completed",
	})

	return &Consumer{
		paymentSuccReader: paymentSuccReader,
		paymentFailReader: paymentFailReader,
		deliveryReader:    deliveryReader,
		service:           service,
	}
}

func (c *Consumer) Close() {
	c.paymentSuccReader.Close()
	c.paymentFailReader.Close()
	c.deliveryReader.Close()
}

func (c *Consumer) consumePaymentSuccess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.paymentSuccReader.ReadMessage(ctx)
			if err != nil {
				log.Println("kafka error:", err)
				continue
			}

			var event models.PaymentSuccessEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json error:", err)
				continue
			}

			if err := c.service.HandlePaymentSuccess(ctx, event); err != nil {
				log.Println("handle error:", err)
			}
		}
	}
}
func (c *Consumer) consumePaymentFailed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.paymentFailReader.ReadMessage(ctx)
			if err != nil {
				log.Println("kafka payment_failed error:", err)
				continue
			}

			var event models.PaymentFailedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json payment_failed error:", err)
				continue
			}

			if err := c.service.HandlePaymentFailed(ctx, event); err != nil {
				log.Println("handle payment_failed error:", err)
			}
		}
	}
}

func (c *Consumer) consumeDelivery(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.deliveryReader.ReadMessage(ctx)
			if err != nil {
				log.Println("kafka delivery_completed error:", err)
				continue
			}

			var event models.DeliveryEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json delivery_completed error:", err)
				continue
			}

			if err := c.service.HandleDelivery(ctx, event); err != nil {
				log.Println("handle delivery_completed error:", err)
			}
		}
	}
}

func (c *Consumer) Run(ctx context.Context) {
	go c.consumePaymentSuccess(ctx)
	go c.consumePaymentFailed(ctx)
	go c.consumeDelivery(ctx)

	<-ctx.Done()
	c.Close()
}
