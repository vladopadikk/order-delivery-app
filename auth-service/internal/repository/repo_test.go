package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/models"
)

func TestRepoCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type args struct {
		username     string
		email        string
		passwordHash string
		createdAt    time.Time
	}

	type mockBehavior func(id int64, args args)

	testCases := []struct {
		name         string
		args         args
		id           int64
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			args: args{
				username:     "petr",
				email:        "petr@mail.com",
				passwordHash: "pass-hash",
				createdAt:    time.Now(),
			},
			id: 42,
			mockBehavior: func(id int64, args args) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						args.username,
						args.email,
						args.passwordHash,
						args.createdAt,
					).WillReturnRows(rows)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.id, tC.args)

			got, err := r.Create(
				context.Background(),
				tC.args.username,
				tC.args.email,
				tC.args.passwordHash,
				tC.args.createdAt,
			)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	r := NewRepository(db)

	type args struct {
		email string
	}

	type mockBehavior func(user models.User, args args)

	testCases := []struct {
		name         string
		args         args
		expectedUser models.User
		wantError    bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			args: args{
				email: "petr@mail.com",
			},
			expectedUser: models.User{
				ID:           42,
				Username:     "petr",
				Email:        "petr@mail.com",
				PasswordHash: "pass-hash",
				CreatedAt:    time.Now(),
			},
			mockBehavior: func(user models.User, args args) {
				rows := sqlmock.NewRows([]string{
					"id",
					"username",
					"email",
					"password_hash",
					"created_at",
				}).
					AddRow(
						user.ID,
						user.Username,
						user.Email,
						user.PasswordHash,
						user.CreatedAt,
					)

				mock.ExpectQuery("SELECT id, username, email, password_hash, created_at FROM users").
					WithArgs(args.email).
					WillReturnRows(rows)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mockBehavior(tC.expectedUser, tC.args)

			got, err := r.GetByEmail(
				context.Background(),
				tC.args.email,
			)

			if tC.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tC.expectedUser, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
