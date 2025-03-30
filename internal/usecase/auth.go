package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo repository.Repository
}

func NewAuthUsecase(repo repository.Repository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) IsEmailTaken(ctx context.Context, email string) bool {
	_, err := u.repo.GetUserByEmail(ctx, email)
	return err == nil
}

func (u *AuthUsecase) CreateUser(ctx context.Context, email, password, firstName, lastName string) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, errors.New("ошибка при хешировании пароля")
	}

	user := repository.User{
		Email:        email,
		Password:     string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		TokenVersion: 1,
	}

	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = id

	return domain.User{
		ID:        int(id),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (u *AuthUsecase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (u *AuthUsecase) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, int64(id))
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
