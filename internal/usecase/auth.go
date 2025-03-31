package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo repository.Repository
	logger logger.Logger
}

func NewAuthUsecase(repo repository.Repository, logger logger.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, logger: logger}
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

	requestID := ctx.Value(utils.RequestIDKey)
	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"err": err.Error(),
		}).Error("User usecase: create user failed")
		return domain.User{}, err
	}
	u.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"user_id": id,
	}).Info("User usecase: user created succesfully")
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
	requestID := ctx.Value(utils.RequestIDKey)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"err": err.Error(),
		}).Error("User usecase: get user by email failed")
		return domain.User{}, err
	}

	u.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"user_id": user.ID,
	}).Info("User usecase: get user by email")

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
	requestID := ctx.Value(utils.RequestIDKey)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"err": err.Error(),
		}).Error("User usecase: get user by id failed")
		return domain.User{}, err
	}

	u.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"user_id": user.ID,
	}).Info("User usecase: get user by id")

	return domain.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
