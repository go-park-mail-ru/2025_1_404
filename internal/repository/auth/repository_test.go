package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (AuthRepository, pgxmock.PgxPoolIface) {
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

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name, u.last_name, u.email, u.password\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.email = \$1`).
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", email, "hashed_pw"))

	u, err := repo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "avatar.png", u.Image)
	require.Equal(t, email, u.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name, u.last_name, u.email, u.password\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", "user@example.com", "hashed_pw"))

	u, err := repo.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "Ivan", u.FirstName)
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
