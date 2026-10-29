package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

const (
	publishTimeout  = 10 * time.Second
	defaultInterval = 500 * time.Millisecond
	defaultBatch    = 50
)

type Repository interface {
	GetUnpublishedOutbox(ctx context.Context, limit int) ([]models.OutboxMessage, error)
	MarkOutboxPublished(ctx context.Context, id int64) error
}

type Publisher interface {
	Publish(ctx context.Context, topic, key string, payload []byte) error
}

// Relay читает неотправленные события из outbox, публикует их в Kafka
// и помечает отправленными. Работает в одном экземпляре на сервис.
type Relay struct {
	repo      Repository
	publisher Publisher
	interval  time.Duration
	batchSize int
}

func NewRelay(repo Repository, publisher Publisher, interval time.Duration, batchSize int) *Relay {
	if interval <= 0 {
		interval = defaultInterval
	}
	if batchSize <= 0 {
		batchSize = defaultBatch
	}

	return &Relay{
		repo:      repo,
		publisher: publisher,
		interval:  interval,
		batchSize: batchSize,
	}
}

// Run крутится до отмены контекста.
func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		n, err := r.processBatch(ctx)
		if err != nil && ctx.Err() == nil {
			log.Println("outbox relay error:", err)
		}

		// Пачка полная - значит, в таблице могут быть ещё события, не ждём тикер.
		if err == nil && n == r.batchSize {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// processBatch отправляет пачку событий по порядку. При первой ошибке останавливается:
// неотправленное событие остаётся в таблице, и следующий цикл начнёт с него же.
func (r *Relay) processBatch(ctx context.Context) (int, error) {
	messages, err := r.repo.GetUnpublishedOutbox(ctx, r.batchSize)
	if err != nil {
		return 0, fmt.Errorf("get unpublished events: %w", err)
	}

	published := 0

	for _, msg := range messages {
		if err := r.publish(ctx, msg); err != nil {
			return published, fmt.Errorf("publish event %s: %w", msg.EventID, err)
		}

		if err := r.repo.MarkOutboxPublished(ctx, msg.ID); err != nil {
			return published, fmt.Errorf("mark event %s as published: %w", msg.EventID, err)
		}

		published++
	}

	return published, nil
}

func (r *Relay) publish(ctx context.Context, msg models.OutboxMessage) error {
	ctx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()

	return r.publisher.Publish(ctx, msg.Topic, msg.Key, msg.Payload)
}
