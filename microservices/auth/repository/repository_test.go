package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*authRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewAuthRepository(mock, logger.NewStub())
	return repo, mock
}

func TestRepository_CreateUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()

	user := User{
		FirstName:    "Ivan",
		LastName:     "Petrov",
		Email:        "ivan@mail.ru",
		Password:     "hashed_pw",
		TokenVersion: 1,
		ImageID:      nil,
	}

	mock.ExpectQuery(`(?i)INSERT INTO kvartirum.Users`).
		WithArgs((*int64)(nil), user.FirstName, user.LastName, user.Email, user.Password, user.TokenVersion).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, int64(42), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByEmail(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	email := "user@example.com"
	id := int64(1)
	createdAt := time.Now()

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name,\s+u.last_name,\s+u.email,\s+u.password,\s+u.created_at,\s+u.role\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.email = \$1`).
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password", "created_at", "role",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", email, "hashed_pw", createdAt, "user"))

	u, err := repo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	require.EqualValues(t, id, u.ID)
	require.Equal(t, "Ivan", u.FirstName)
	require.Equal(t, "avatar.png", u.Image)
	require.Equal(t, createdAt, u.CreatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)
	createdAt := time.Now()

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name,\s+u.last_name,\s+u.email,\s+u.password,\s+u.created_at,\s+u.role\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password", "created_at", "role",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", "user@example.com", "hashed_pw", createdAt, "user"))

	u, err := repo.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.EqualValues(t, id, u.ID)
	require.Equal(t, "Ivan", u.FirstName)
	require.Equal(t, createdAt, u.CreatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	user := domain.User{
		ID:        1,
		FirstName: "Ivan",
		LastName:  "Petrov",
		Email:     "new@mail.ru",
		Image:     "new.png",
	}

	mock.ExpectQuery(`(?i)UPDATE kvartirum.Users`).
		WithArgs(user.Image, user.FirstName, user.LastName, user.Email, user.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "first_name", "last_name", "email", "image_uuid",
		}).AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Image))

	updated, err := repo.UpdateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, user.ID, updated.ID)
	require.Equal(t, user.Image, updated.Image)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.Users WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteUser(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateImage(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	filename := "image-uuid.png"

	mock.ExpectExec(`(?i)INSERT INTO kvartirum.Image \(uuid\) VALUES \(\$1\)`).
		WithArgs(filename).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateImage(ctx, filename)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetImageByUUID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	uuid := "image-uuid"
	expectedID := int64(123)

	mock.ExpectQuery(`(?i)SELECT id FROM kvartirum.Image WHERE uuid = \$1`).
		WithArgs(uuid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := repo.GetImageByUUID(ctx, uuid)
	require.NoError(t, err)
	require.True(t, id.Valid)
	require.Equal(t, expectedID, id.Int64)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetImageByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	id := sql.NullInt64{Int64: 123, Valid: true}
	expectedUUID := "image-uuid.png"

	mock.ExpectQuery(`(?i)SELECT uuid from kvartirum.Image WHERE id = \$1`).
		WithArgs(id.Int64).
		WillReturnRows(pgxmock.NewRows([]string{"uuid"}).AddRow(expectedUUID))

	uuid, err := repo.GetImageByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expectedUUID, uuid)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteUserImage(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	userID := int64(1)

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.Image where id = \(\s*SELECT image_id from kvartirum.Users where id = \$1\s*\)`).
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteUserImage(ctx, userID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
