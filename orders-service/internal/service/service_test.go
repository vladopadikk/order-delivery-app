package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	dmocks "github.com/vladopadikk/order-delivery-app/orders-service/internal/database/mocks"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/service/mocks"
)

func TestCreateOrderService(t *testing.T) {
	type mockBehavior func(
		m *mocks.MockOrderRepository,
		tx *dmocks.MockTx, prod *mocks.MockProducer,
		userID int64,
		in models.OrderInput,
	)

	testCases := []struct {
		name             string
		userID           int64
		orderIn          models.OrderInput
		mockBehavior     mockBehavior
		expectedResponse models.OrderResponse
		wantErr          bool
	}{
		{
			name:   "OK",
			userID: 42,
			orderIn: models.OrderInput{
				Items:           []models.OrderItemInput{{ProductID: 1, Quantity: 2, Price: 500}},
				DeliveryAddress: "Moscow, Lenina 1",
			},
			mockBehavior: func(
				m *mocks.MockOrderRepository,
				tx *dmocks.MockTx,
				prod *mocks.MockProducer,
				userID int64,
				in models.OrderInput,
			) {
				m.EXPECT().BeginTx(mock.Anything).Return(tx, nil)

				m.EXPECT().Create(
					mock.Anything,
					tx,
					userID,
					models.StatusPaymentPending,
					1000.0,
					in.DeliveryAddress,
				).Return(models.Order{
					ID:              52,
					UserID:          userID,
					Status:          models.StatusPaymentPending,
					TotalPrice:      1000,
					DeliveryAddress: in.DeliveryAddress,
				}, nil)

				m.EXPECT().
					CreateItems(mock.Anything, tx, int64(52), in.Items[0]).
					Return(nil)

				tx.EXPECT().Commit().Return(nil)
				tx.EXPECT().Rollback().Return(nil).Maybe()

				prod.EXPECT().
					PublishOrderCreated(mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedResponse: models.OrderResponse{
				OrderID:    52,
				TotalPrice: 1000,
				Status:     models.StatusPaymentPending,
			},
		},
		{
			name:    "BeginTx error",
			userID:  1,
			orderIn: models.OrderInput{Items: []models.OrderItemInput{{ProductID: 10, Quantity: 2, Price: 500}}},
			mockBehavior: func(
				m *mocks.MockOrderRepository,
				tx *dmocks.MockTx,
				prod *mocks.MockProducer,
				userID int64,
				in models.OrderInput,
			) {
				m.EXPECT().BeginTx(mock.Anything).Return(nil, errors.New("db down"))
			},
			wantErr: true,
		},
		{
			name:    "Create error",
			userID:  1,
			orderIn: models.OrderInput{Items: []models.OrderItemInput{{ProductID: 10, Quantity: 2, Price: 500}}},
			mockBehavior: func(
				m *mocks.MockOrderRepository,
				tx *dmocks.MockTx,
				prod *mocks.MockProducer,
				userID int64,
				in models.OrderInput,
			) {
				m.EXPECT().BeginTx(mock.Anything).Return(tx, nil)

				m.EXPECT().Create(
					mock.Anything,
					tx,
					userID,
					models.StatusPaymentPending,
					1000.0,
					in.DeliveryAddress,
				).Return(models.Order{}, errors.New("insert failed"))

				tx.EXPECT().Rollback().Return(nil)
			},
			wantErr: true,
		},
		{
			name:    "CreateItems error",
			userID:  1,
			orderIn: models.OrderInput{Items: []models.OrderItemInput{{ProductID: 10, Quantity: 2, Price: 500}}},
			mockBehavior: func(
				m *mocks.MockOrderRepository,
				tx *dmocks.MockTx,
				prod *mocks.MockProducer,
				userID int64,
				in models.OrderInput,
			) {
				m.EXPECT().BeginTx(mock.Anything).Return(tx, nil)

				m.EXPECT().Create(
					mock.Anything,
					tx,
					userID,
					models.StatusPaymentPending,
					1000.0,
					in.DeliveryAddress,
				).Return(models.Order{
					ID:              52,
					UserID:          userID,
					Status:          models.StatusPaymentPending,
					TotalPrice:      1000,
					DeliveryAddress: in.DeliveryAddress,
				}, nil)

				m.EXPECT().
					CreateItems(mock.Anything, tx, int64(52), in.Items[0]).
					Return(errors.New("insert failed"))

				tx.EXPECT().Rollback().Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			repo := mocks.NewMockOrderRepository(t)
			tx := dmocks.NewMockTx(t)
			prod := mocks.NewMockProducer(t)

			tC.mockBehavior(repo, tx, prod, tC.userID, tC.orderIn)

			svc := NewService(repo, prod)

			got, err := svc.CreateOrder(
				context.Background(),
				tC.userID,
				tC.orderIn,
			)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.expectedResponse, got)
			}
		})
	}
}
