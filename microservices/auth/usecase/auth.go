package usecase

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo   auth.AuthRepository
	logger logger.Logger
	s3Repo s3.S3Repo
	cfg    *config.Config
}

func NewAuthUsecase(repo auth.AuthRepository, logger logger.Logger, s3Repo s3.S3Repo, cfg *config.Config) *authUsecase {
	return &authUsecase{repo: repo, logger: logger, s3Repo: s3Repo, cfg: cfg}
}

func (u *authUsecase) IsEmailTaken(ctx context.Context, email string) bool {
	_, err := u.repo.GetUserByEmail(ctx, email)
	return err == nil
}

func (u *authUsecase) CreateUser(ctx context.Context, email, password, firstName, lastName string) (domain.User, error) {
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

	if u.IsEmailTaken(ctx, email) {
		return domain.User{}, errors.New("email уже занят")
	}

	requestID := ctx.Value(utils.RequestIDKey)
	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("User usecase: create user failed")
		return domain.User{}, err
	}

	user.ID = id

	return domain.User{
		ID:        int(id),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func (u *authUsecase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("User usecase: get user by email failed")
		return domain.User{}, err
	}

	if user.Image != "" {
		user.Image = u.cfg.Minio.Path + u.cfg.Minio.AvatarsBucket + user.Image
	}

	return user, nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	user, err := u.repo.GetUserByID(ctx, int64(id))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("User usecase: get user by id failed")
		return domain.User{}, err
	}

	if user.Image != "" {
		user.Image = u.cfg.Minio.Path + u.cfg.Minio.AvatarsBucket + user.Image
	}

	return user, nil
}

func (u *authUsecase) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	currentUser, err := u.GetUserByID(ctx, user.ID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": user.ID}).Warn("user id not found")
		return domain.User{}, errors.New("пользователь не найден")
	}

	if currentUser.Email != user.Email {
		if u.IsEmailTaken(ctx, user.Email) {
			return domain.User{}, errors.New("email уже занят")
		}
	}

	if currentUser.Image != "" {
		currentUser.Image = path.Base(currentUser.Image)
	}

	// Заполняем непереданные поля уже имеющимися
	if user.Email != "" {
		currentUser.Email = user.Email
	}
	if user.FirstName != "" {
		currentUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		currentUser.LastName = user.LastName
	}
	if user.Image != "" {
		currentUser.Image = user.Image
	}

	// Обновляем в БД
	updatedUser, err := u.repo.UpdateUser(ctx, currentUser)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()})
		return domain.User{}, err
	}

	if updatedUser.Image != "" {
		updatedUser.Image = u.cfg.Minio.Path + u.cfg.Minio.AvatarsBucket + updatedUser.Image
	}

	return updatedUser, nil
}

func (u *authUsecase) UploadImage(ctx context.Context, id int, file s3.Upload) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	user, err := u.GetUserByID(ctx, id)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Warn("user id not found")
		return domain.User{}, fmt.Errorf("failed to find user")
	}

	previousImage := user.Image

	// Загружаем в файловое хранилище фото
	fileName, err := u.s3Repo.Put(ctx, file)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("upload image failed")
		return domain.User{}, err
	}

	// Создаем запись в БД
	err = u.repo.CreateImage(ctx, fileName)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("failed to load user image")
		return domain.User{}, err
	}

	if previousImage != "" {
		_, err := u.DeleteImage(ctx, id)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("failed to delete old image")
			return domain.User{}, err
		}
	}

	// Обновляем пользователя с новым именем аватарки
	updatedUser, err := u.UpdateUser(ctx, domain.User{
		ID:    user.ID,
		Image: fileName,
	})

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("failed to update user")
		return domain.User{}, err
	}

	return updatedUser, nil
}

func (u *authUsecase) DeleteImage(ctx context.Context, id int) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	user, err := u.GetUserByID(ctx, id)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Warn("user id not found")
		return domain.User{}, fmt.Errorf("failed to find user")
	}

	err = u.s3Repo.Remove(ctx, "avatars", path.Base(user.Image))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("delete image failed")
		return domain.User{}, err
	}

	err = u.repo.DeleteUserImage(ctx, int64(id))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("failed to delete user image")
		return domain.User{}, err
	}

	updatedUser, err := u.UpdateUser(ctx, domain.User{
		ID:    id,
		Image: "",
	})

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("failed to update user")
		return domain.User{}, err
	}

	return updatedUser, nil
}
