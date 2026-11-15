package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

func TestRepoCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type args struct {
		userID          int64
		status          string
		totalPrice      float64
		deliveryAddress string
	}

	type mockBehavior func(order models.Order, args args)

	testCases := []struct {
		name          string
		args          args
		wantErr       bool
		mockBehavior  mockBehavior
		expectedOrder models.Order
	}{
		{
			name: "OK",
			args: args{
				userID:          42,
				status:          models.StatusPaymentPending,
				totalPrice:      1000.0,
				deliveryAddress: "Moscow, Lenina 1",
			},
			mockBehavior: func(order models.Order, args args) {
				rows := sqlmock.NewRows([]string{
					"id",
					"user_id",
					"status",
					"total_price",
					"delivery_address",
					"created_at",
				}).AddRow(
					order.ID,
					order.UserID,
					order.Status,
					order.TotalPrice,
					order.DeliveryAddress,
					order.CreatedAt,
				)

				mock.ExpectQuery("INSERT INTO orders").
					WithArgs(
						args.userID,
						args.status,
						args.totalPrice,
						args.deliveryAddress,
					).WillReturnRows(rows)
			},
			expectedOrder: models.Order{
				ID:              52,
				UserID:          42,
				Status:          models.StatusPaymentPending,
				TotalPrice:      1000.0,
				DeliveryAddress: "Moscow, Lenina 1",
				CreatedAt:       time.Time{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.expectedOrder, tC.args)

			got, err := r.Create(
				context.Background(),
				db,
				tC.args.userID,
				tC.args.status,
				tC.args.totalPrice,
				tC.args.deliveryAddress,
			)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.expectedOrder, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type mockBehavior func(userID int64, orders []models.OrderResponse)

	testCases := []struct {
		name           string
		userID         int64
		wantErr        bool
		mockBehavior   mockBehavior
		expectedOrders []models.OrderResponse
	}{
		{
			name:   "OK",
			userID: 42,
			mockBehavior: func(userID int64, orders []models.OrderResponse) {
				rows := sqlmock.NewRows([]string{"id", "total_price", "status"})
				for _, o := range orders {
					rows.AddRow(o.OrderID, o.TotalPrice, o.Status)
				}

				mock.ExpectQuery("SELECT o.id, o.total_price, o.status").
					WithArgs(userID).
					WillReturnRows(rows)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.userID, tC.expectedOrders)

			got, err := r.GetOrders(context.Background(), tC.userID)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.expectedOrders, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepoCreateOutboxMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type mockBehavior func(msg models.OutboxMessage)

	testCases := []struct {
		name         string
		msg          models.OutboxMessage
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			msg: models.OutboxMessage{
				EventID:     "7f1c2d3e-0000-4000-8000-000000000001",
				AggregateID: 52,
				Topic:       models.TopicOrderCreated,
				Key:         "52",
				Payload:     []byte(`{"order_id":52}`),
			},
			mockBehavior: func(msg models.OutboxMessage) {
				mock.ExpectExec("INSERT INTO outbox").
					WithArgs(
						msg.EventID,
						msg.AggregateID,
						msg.Topic,
						msg.Key,
						string(msg.Payload),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Insert error",
			msg: models.OutboxMessage{
				EventID:     "7f1c2d3e-0000-4000-8000-000000000002",
				AggregateID: 53,
				Topic:       models.TopicOrderCreated,
				Key:         "53",
				Payload:     []byte(`{"order_id":53}`),
			},
			wantErr: true,
			mockBehavior: func(msg models.OutboxMessage) {
				mock.ExpectExec("INSERT INTO outbox").
					WithArgs(
						msg.EventID,
						msg.AggregateID,
						msg.Topic,
						msg.Key,
						string(msg.Payload),
					).
					WillReturnError(errors.New("insert failed"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.msg)

			err := r.CreateOutboxMessage(context.Background(), db, tC.msg)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUnpublishedOutbox(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type mockBehavior func(limit int, messages []models.OutboxMessage)

	createdAt := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name             string
		limit            int
		wantErr          bool
		mockBehavior     mockBehavior
		expectedMessages []models.OutboxMessage
	}{
		{
			name:  "OK",
			limit: 10,
			mockBehavior: func(limit int, messages []models.OutboxMessage) {
				rows := sqlmock.NewRows([]string{
					"id",
					"event_id",
					"aggregate_id",
					"topic",
					"message_key",
					"payload",
					"created_at",
				})
				for _, m := range messages {
					rows.AddRow(
						m.ID,
						m.EventID,
						m.AggregateID,
						m.Topic,
						m.Key,
						string(m.Payload),
						m.CreatedAt,
					)
				}

				mock.ExpectQuery("FROM outbox").
					WithArgs(limit).
					WillReturnRows(rows)
			},
			expectedMessages: []models.OutboxMessage{
				{
					ID:          1,
					EventID:     "7f1c2d3e-0000-4000-8000-000000000001",
					AggregateID: 52,
					Topic:       models.TopicOrderCreated,
					Key:         "52",
					Payload:     []byte(`{"order_id":52}`),
					CreatedAt:   createdAt,
				},
				{
					ID:          2,
					EventID:     "7f1c2d3e-0000-4000-8000-000000000002",
					AggregateID: 53,
					Topic:       models.TopicOrderCreated,
					Key:         "53",
					Payload:     []byte(`{"order_id":53}`),
					CreatedAt:   createdAt,
				},
			},
		},
		{
			name:  "No unpublished events",
			limit: 10,
			mockBehavior: func(limit int, messages []models.OutboxMessage) {
				rows := sqlmock.NewRows([]string{
					"id",
					"event_id",
					"aggregate_id",
					"topic",
					"message_key",
					"payload",
					"created_at",
				})

				mock.ExpectQuery("FROM outbox").
					WithArgs(limit).
					WillReturnRows(rows)
			},
		},
		{
			name:    "Query error",
			limit:   10,
			wantErr: true,
			mockBehavior: func(limit int, messages []models.OutboxMessage) {
				mock.ExpectQuery("FROM outbox").
					WithArgs(limit).
					WillReturnError(errors.New("db down"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.limit, tC.expectedMessages)

			got, err := r.GetUnpublishedOutbox(context.Background(), tC.limit)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.expectedMessages, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMarkOutboxPublished(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type mockBehavior func(id int64)

	testCases := []struct {
		name         string
		id           int64
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(id int64) {
				mock.ExpectExec("UPDATE outbox").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:    "Update error",
			id:      2,
			wantErr: true,
			mockBehavior: func(id int64) {
				mock.ExpectExec("UPDATE outbox").
					WithArgs(id).
					WillReturnError(errors.New("db down"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.id)

			err := r.MarkOutboxPublished(context.Background(), tC.id)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
