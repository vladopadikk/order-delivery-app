package outbox

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/outbox/mocks"
)

func TestProcessBatch(t *testing.T) {
	type mockBehavior func(
		repo *mocks.MockRepository,
		pub *mocks.MockPublisher,
		messages []models.OutboxMessage,
	)

	messages := []models.OutboxMessage{
		{ID: 1, EventID: "e1", Topic: models.TopicOrderCreated, Key: "1", Payload: []byte(`{"order_id":1}`)},
		{ID: 2, EventID: "e2", Topic: models.TopicOrderCreated, Key: "2", Payload: []byte(`{"order_id":2}`)},
		{ID: 3, EventID: "e3", Topic: models.TopicOrderCreated, Key: "3", Payload: []byte(`{"order_id":3}`)},
	}

	testCases := []struct {
		name              string
		messages          []models.OutboxMessage
		mockBehavior      mockBehavior
		expectedPublished int
		wantErr           bool
	}{
		{
			name:     "OK",
			messages: messages,
			mockBehavior: func(
				repo *mocks.MockRepository,
				pub *mocks.MockPublisher,
				messages []models.OutboxMessage,
			) {
				repo.EXPECT().GetUnpublishedOutbox(mock.Anything, 10).Return(messages, nil)

				for _, msg := range messages {
					pub.EXPECT().
						Publish(mock.Anything, msg.Topic, msg.Key, msg.Payload).
						Return(nil)

					repo.EXPECT().
						MarkOutboxPublished(mock.Anything, msg.ID).
						Return(nil)
				}
			},
			expectedPublished: 3,
		},
		{
			name:     "No unpublished events",
			messages: nil,
			mockBehavior: func(
				repo *mocks.MockRepository,
				pub *mocks.MockPublisher,
				messages []models.OutboxMessage,
			) {
				repo.EXPECT().GetUnpublishedOutbox(mock.Anything, 10).Return(nil, nil)
			},
			expectedPublished: 0,
		},
		{
			name:     "GetUnpublishedOutbox error",
			messages: messages,
			mockBehavior: func(
				repo *mocks.MockRepository,
				pub *mocks.MockPublisher,
				messages []models.OutboxMessage,
			) {
				repo.EXPECT().
					GetUnpublishedOutbox(mock.Anything, 10).
					Return(nil, errors.New("db down"))
			},
			wantErr: true,
		},
		{
			// Ошибка на втором событии: первое отправлено и помечено, второе остаётся
			// неотправленным, третье не трогаем, чтобы не нарушить порядок.
			name:     "Publish error stops the batch",
			messages: messages,
			mockBehavior: func(
				repo *mocks.MockRepository,
				pub *mocks.MockPublisher,
				messages []models.OutboxMessage,
			) {
				repo.EXPECT().GetUnpublishedOutbox(mock.Anything, 10).Return(messages, nil)

				pub.EXPECT().
					Publish(mock.Anything, messages[0].Topic, messages[0].Key, messages[0].Payload).
					Return(nil)
				repo.EXPECT().
					MarkOutboxPublished(mock.Anything, messages[0].ID).
					Return(nil)

				pub.EXPECT().
					Publish(mock.Anything, messages[1].Topic, messages[1].Key, messages[1].Payload).
					Return(errors.New("kafka unavailable"))
			},
			expectedPublished: 1,
			wantErr:           true,
		},
		{
			name:     "MarkOutboxPublished error",
			messages: messages,
			mockBehavior: func(
				repo *mocks.MockRepository,
				pub *mocks.MockPublisher,
				messages []models.OutboxMessage,
			) {
				repo.EXPECT().GetUnpublishedOutbox(mock.Anything, 10).Return(messages, nil)

				pub.EXPECT().
					Publish(mock.Anything, messages[0].Topic, messages[0].Key, messages[0].Payload).
					Return(nil)
				repo.EXPECT().
					MarkOutboxPublished(mock.Anything, messages[0].ID).
					Return(errors.New("db down"))
			},
			expectedPublished: 0,
			wantErr:           true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			pub := mocks.NewMockPublisher(t)

			tC.mockBehavior(repo, pub, tC.messages)

			relay := NewRelay(repo, pub, time.Second, 10)

			got, err := relay.processBatch(context.Background())

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tC.expectedPublished, got)
		})
	}
}

func TestRunRetriesAfterPublishError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	pub := mocks.NewMockPublisher(t)

	msg := models.OutboxMessage{
		ID:      1,
		EventID: "e1",
		Topic:   models.TopicOrderCreated,
		Key:     "1",
		Payload: []byte(`{"order_id":1}`),
	}

	var marked atomic.Bool

	repo.EXPECT().
		GetUnpublishedOutbox(mock.Anything, 10).
		RunAndReturn(func(ctx context.Context, limit int) ([]models.OutboxMessage, error) {
			if marked.Load() {
				return nil, nil
			}
			return []models.OutboxMessage{msg}, nil
		})

	repo.EXPECT().
		MarkOutboxPublished(mock.Anything, msg.ID).
		RunAndReturn(func(ctx context.Context, id int64) error {
			marked.Store(true)
			return nil
		}).
		Once()

	// Первая попытка падает, вторая проходит.
	pub.EXPECT().
		Publish(mock.Anything, msg.Topic, msg.Key, msg.Payload).
		Return(errors.New("kafka unavailable")).
		Once()
	pub.EXPECT().
		Publish(mock.Anything, msg.Topic, msg.Key, msg.Payload).
		Return(nil).
		Once()

	relay := NewRelay(repo, pub, 5*time.Millisecond, 10)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		relay.Run(ctx)
		close(done)
	}()

	assert.Eventually(t, marked.Load, 2*time.Second, 5*time.Millisecond)

	cancel()
	<-done
}
