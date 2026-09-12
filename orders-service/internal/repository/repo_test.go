package repository

import (
	"context"
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
