package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
)

//go:generate mockgen -source interface.go -destination=mocks/mock_auth_repo.go -package=mocks

type AuthRepository interface {
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateImage(ctx context.Context, file filestorage.FileUpload) error
	GetImageByID(ctx context.Context, id sql.NullInt64) (string, error)
	DeleteUserImage(ctx context.Context, id int64) error
}
