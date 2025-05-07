package usecase

import (
	"context"
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	mock "github.com/go-park-mail-ru/2025_1_404/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	mockS3 "github.com/go-park-mail-ru/2025_1_404/pkg/database/s3/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	Path   = "http://localhost:9000"
	Bucket = "/avatars/"
)

func TestUploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	file := s3.Upload{
		Filename: "avatar.png",
		Bucket:   "avatar",
	}

	user := domain.User{
		ID:    1,
		Image: "image.png",
	}

	t.Run("UploadImage ok", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(domain.User{ID: 1}, nil).Times(2)
		mockS3.EXPECT().Put(ctx, file).Return(file.Filename, nil)
		mockRepo.EXPECT().CreateImage(ctx, file.Filename).Return(nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any(), domain.User{ID: 1, Image: file.Filename}).Return(user, nil)

		updatedUser, err := usecase.UploadImage(ctx, 1, file)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, updatedUser.ID)
		assert.Equal(t, Path+Bucket+user.Image, updatedUser.Image)
	})

	t.Run("GetUserById faield", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("not found"))

		_, err := usecase.UploadImage(ctx, 1, file)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find user")
	})

	t.Run("Put image failed", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(domain.User{ID: 1}, nil)
		mockS3.EXPECT().Put(ctx, file).Return("", errors.New("put error"))

		_, err := usecase.UploadImage(ctx, 1, file)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "put error")
	})

	t.Run("CreateImage failed", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(domain.User{ID: 1}, nil)
		mockS3.EXPECT().Put(ctx, file).Return(file.Filename, nil)
		mockRepo.EXPECT().CreateImage(ctx, file.Filename).Return(errors.New("create image error"))

		_, err := usecase.UploadImage(ctx, 1, file)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create image error")
	})

}

func TestIsEmailTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	user := domain.User{
		ID:    1,
		Email: "e@mail.ru",
	}
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Return(user, nil)

	assert.Equal(t, true, usecase.IsEmailTaken(ctx, user.Email))
}

func TestCreateUser(t *testing.T) {
	email := "e@mail.ru"
	password := "password"
	firstName := "Maksim"
	lastName := "Maksimov"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(domain.User{}, errors.New("user not found"))
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(int64(1), nil)

		user, err := usecase.CreateUser(ctx, email, password, firstName, lastName)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)
		assert.NotEmpty(t, user.Password)
	})

	t.Run("email is taken", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(domain.User{ID: 123}, nil)

		_, err := usecase.CreateUser(ctx, email, password, firstName, lastName)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email уже занят")
	})

	t.Run("repo.CreateUser failed", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(domain.User{}, errors.New("user not found"))
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))

		_, err := usecase.CreateUser(ctx, email, password, firstName, lastName)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	t.Run("GetUserByEmail ok", func(t *testing.T) {
		expectedUser := domain.User{
			Email: "e@mail.ru",
			Image: Path + Bucket + "avatar.png",
		}

		user := domain.User{
			Email: "e@mail.ru",
			Image: "avatar.png",
		}
		mockRepo.EXPECT().GetUserByEmail(gomock.Any(), expectedUser.Email).Return(user, nil)

		user, err := usecase.GetUserByEmail(ctx, expectedUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, user, expectedUser)
	})

	t.Run("GetUserByEmail repository failed", func(t *testing.T) {
		email := "e@mail.ru"
		mockRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Return(domain.User{}, fmt.Errorf("repository GetUserByEmail failed"))

		user, err := usecase.GetUserByEmail(ctx, email)

		assert.Error(t, err)
		assert.EqualError(t, err, "repository GetUserByEmail failed")
		assert.Equal(t, user, domain.User{})
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	t.Run("UpdateUser ok", func(t *testing.T) {
		user := domain.User{
			ID:        1,
			Email:     "new@mail.ru",
			FirstName: "Ivan",
			LastName:  "Ivanov",
		}

		currentUser := domain.User{
			ID:    1,
			Email: "old@mail.ru",
		}

		updatedUser := user

		mockRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(domain.User{}, fmt.Errorf("user not found"))
		mockRepo.EXPECT().GetUserByID(ctx, int64(1)).Return(currentUser, nil)
		mockRepo.EXPECT().UpdateUser(ctx, user).Return(updatedUser, nil)

		updatedUser, err := usecase.UpdateUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, user, updatedUser)
	})

	t.Run("GetUserByID repository failed", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, int64(999)).Return(domain.User{}, fmt.Errorf("GetUserByID repository failed"))

		updatedUser, err := usecase.UpdateUser(ctx, domain.User{ID: 999})
		assert.Error(t, err)
		assert.Equal(t, updatedUser, domain.User{})
	})

	t.Run("UpdateUser repository failed", func(t *testing.T) {
		user := domain.User{ID: 1}
		mockRepo.EXPECT().GetUserByID(ctx, int64(1)).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(ctx, user).Return(domain.User{}, fmt.Errorf("UpdateUser repository failed"))

		updatedUser, err := usecase.UpdateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, updatedUser, domain.User{})
	})
}

func TestDeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuthRepository(ctrl)
	mockS3 := mockS3.NewMockS3Repo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{Minio: config.MinioConfig{Path: Path, AvatarsBucket: "/avatars/"}}

	usecase := NewAuthUsecase(mockRepo, logger, mockS3, cfg)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	userID := 1
	imagePath := "path/to/image.jpg"

	t.Run("DeleteImage ok", func(t *testing.T) {
		existingUser := domain.User{
			ID:    userID,
			Image: imagePath,
		}
		updatedUser := domain.User{
			ID:    userID,
			Image: "",
		}

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(existingUser, nil).Times(2)
		mockS3.EXPECT().Remove(ctx, "avatars", path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(nil)
		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(updatedUser, nil)

		result, err := usecase.DeleteImage(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, updatedUser, result)
	})

	t.Run("user not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("user not found")

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(domain.User{}, expectedErr)

		result, err := usecase.DeleteImage(ctx, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to find user")
		assert.Equal(t, domain.User{}, result)
	})

	t.Run("file storage deletion failed", func(t *testing.T) {
		existingUser := domain.User{
			ID:    userID,
			Image: imagePath,
		}
		fsErr := fmt.Errorf("storage error")

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(existingUser, nil)
		mockS3.EXPECT().Remove(ctx, "avatars", path.Base(imagePath)).Return(fsErr)

		result, err := usecase.DeleteImage(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, fsErr, err)
		assert.Equal(t, domain.User{}, result)
	})

	t.Run("repository DeleteUserImage failed", func(t *testing.T) {
		existingUser := domain.User{
			ID:    userID,
			Image: imagePath,
		}
		repoErr := fmt.Errorf("repository error")

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(existingUser, nil)
		mockS3.EXPECT().Remove(ctx, "avatars", path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(repoErr)

		result, err := usecase.DeleteImage(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, domain.User{}, result)
	})

	t.Run("user update failed", func(t *testing.T) {
		existingUser := domain.User{
			ID:    userID,
			Image: imagePath,
		}
		updateErr := fmt.Errorf("update error")

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(existingUser, nil).Times(2)
		mockS3.EXPECT().Remove(ctx, "avatars", path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(nil)
		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(domain.User{}, updateErr)

		result, err := usecase.DeleteImage(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, updateErr, err)
		assert.Equal(t, domain.User{}, result)
	})
}
