package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_auth.go -package=mocks

type authUsecase interface {
	IsEmailTaken(ctx context.Context, email string) bool
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	UploadImage(ctx context.Context, id int, file filestorage.FileUpload) (domain.User, error)
	DeleteImage(ctx context.Context, id int) (domain.User, error)
}
