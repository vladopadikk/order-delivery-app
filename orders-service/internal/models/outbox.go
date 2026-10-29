package models

import (
	"crypto/rand"
	"fmt"
	"time"
)

// OutboxMessage - строка таблицы outbox: событие, которое нужно отправить в Kafka.
type OutboxMessage struct {
	ID          int64
	EventID     string
	AggregateID int64
	Topic       string
	Key         string
	Payload     []byte
	CreatedAt   time.Time
}

// NewEventID возвращает случайный идентификатор события (UUID v4).
func NewEventID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate event id: %w", err)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
