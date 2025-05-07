package auth

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/repository"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_auth_repo.go -package=mocks

type AuthRepository interface {
	CreateUser(ctx context.Context, user repository.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateImage(ctx context.Context, fileName string) error
	GetImageByID(ctx context.Context, id sql.NullInt64) (string, error)
	DeleteUserImage(ctx context.Context, id int64) error
}
