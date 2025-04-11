package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	l, _ := logger.NewZapLogger()

	repo := NewRepository(mock, l)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, 1)

	testEmail := "email@mail.ru"
	expectedUser := User{
		ID:                 1,
		ImageID:            nil,
		FirstName:          "Иван",
		LastName:           "Смирнов",
		Email:              testEmail,
		Password:           "hashed_password",
		LastNotificationID: nil,
		TokenVersion:       1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	mock.ExpectQuery(`SELECT id, image_id, first_name, last_name, email, password, last_notification_id, token_version, created_at, updated_at FROM kvartirum.Users WHERE email = \$1`).
		WithArgs(testEmail).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image_id", "first_name", "last_name", "email", "password", "last_notification_id", "token_version", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID,
			expectedUser.ImageID,
			expectedUser.FirstName,
			expectedUser.LastName,
			expectedUser.Email,
			expectedUser.Password,
			expectedUser.LastNotificationID,
			expectedUser.TokenVersion,
			expectedUser.CreatedAt,
			expectedUser.UpdatedAt,
		))

	user, err := repo.GetUserByEmail(ctx, testEmail)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

