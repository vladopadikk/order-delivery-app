package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

const (
	OrderCreatedTopic      = "order_created"
	PaymentSuccessTopic    = "payment_success"
	PaymentFailedTopic     = "payment_failed"
	DeliveryCompletedTopic = "delivery_completed"
)

func InitTopics(broker string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerConn, err := kafka.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port),
	)
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{
			Topic:             OrderCreatedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             PaymentSuccessTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             PaymentFailedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             DeliveryCompletedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		return fmt.Errorf("create topics: %w", err)
	}

	return nil
}
