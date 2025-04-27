package usecase

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	mockFS "github.com/go-park-mail-ru/2025_1_404/internal/filestorage/mocks"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository/auth"
	mockRepo "github.com/go-park-mail-ru/2025_1_404/internal/usecase/auth/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger := logger.NewStub()

	usecase := NewAuthUsecase(mockRepo, logger, mockFS)

	file := filestorage.FileUpload{
		Name:        "avatar.png",
		Size:        1024,
		File:        nil,
		ContentType: "png",
	}

	user := domain.User{ID: 1, Image: file.Name}

	mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(domain.User{ID: 1}, nil).Times(2)
	mockFS.EXPECT().Add(file).Return(nil)
	mockRepo.EXPECT().CreateImage(gomock.Any(), file).Return(nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any(), domain.User{ID: 1, Image: file.Name}).Return(user, nil)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "333")
	updatedUser, err := usecase.UploadImage(ctx, 1, file)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, updatedUser.ID)
	assert.Equal(t, utils.BasePath+utils.ImagesPath+user.Image, updatedUser.Image)
}

func TestIsEmailTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger := logger.NewStub()

	usecaseAuth := NewAuthUsecase(mockRepo, logger, mockFS)

	user := domain.User{
		ID:    1,
		Email: "e@mail.ru",
	}
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Return(user, nil)

	assert.Equal(t, true, usecaseAuth.IsEmailTaken(context.Background(), user.Email))
}

func TestCreateUser(t *testing.T) {
	email := "e@mail.ru"
	password := "password"
	firstName := "Maksim"
	lastName := "Maksimov"
	requestID := "1"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger := logger.NewStub()

	usecaseAuth := NewAuthUsecase(mockRepo, logger, mockFS)

	t.Run("CreateUser ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.RequestIDKey, requestID)

		mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(domain.User{}, fmt.Errorf("user not found"))
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, u repository.User) (int64, error) {
				err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
				assert.NoError(t, err)

				assert.Equal(t, email, u.Email)
				assert.Equal(t, firstName, u.FirstName)
				assert.Equal(t, lastName, u.LastName)
				assert.Equal(t, 1, u.TokenVersion)

				return int64(1), nil
			})

		user, err := usecaseAuth.CreateUser(ctx, email, password, firstName, lastName)
		assert.NoError(t, err)

		assert.Equal(t, 1, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		assert.NoError(t, err)
	})
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger := logger.NewStub()

	usecaseAuth := NewAuthUsecase(mockRepo, logger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "1")
	t.Run("GetUserByEmail ok", func(t *testing.T) {
		expectedUser := domain.User{
			Email: "e@mail.ru",
			Image: utils.BasePath + utils.ImagesPath + "avatar.png",
		}

		user := domain.User{
			Email: "e@mail.ru",
			Image: "avatar.png",
		}
		mockRepo.EXPECT().GetUserByEmail(gomock.Any(), expectedUser.Email).Return(user, nil)

		user, err := usecaseAuth.GetUserByEmail(ctx, expectedUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, user, expectedUser)
	})

	t.Run("GetUserByEmail repository failed", func(t *testing.T) {
		email := "e@mail.ru"
		mockRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Return(domain.User{}, fmt.Errorf("repository GetUserByEmail failed"))

		user, err := usecaseAuth.GetUserByEmail(ctx, email)

		assert.Error(t, err)
		assert.EqualError(t, err, "repository GetUserByEmail failed")
		assert.Equal(t, user, domain.User{})
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger := logger.NewStub()

	usecaseAuth := NewAuthUsecase(mockRepo, logger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "1")
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

		updatedUser, err := usecaseAuth.UpdateUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, user, updatedUser)
	})

	t.Run("GetUserByID repository failed", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, int64(999)).Return(domain.User{}, fmt.Errorf("GetUserByID repository failed"))

		updatedUser, err := usecaseAuth.UpdateUser(ctx, domain.User{ID: 999})
		assert.Error(t, err)
		assert.Equal(t, updatedUser, domain.User{})
	})

	t.Run("UpdateUser repository failed", func(t *testing.T) {
		user := domain.User{ID: 1}
		mockRepo.EXPECT().GetUserByID(ctx, int64(1)).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(ctx, user).Return(domain.User{}, fmt.Errorf("UpdateUser repository failed"))

		updatedUser, err := usecaseAuth.UpdateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, updatedUser, domain.User{})
	})
}

func TestDeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockauthRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	mockLogger := logger.NewStub()

	usecaseAuth := NewAuthUsecase(mockRepo, mockLogger, mockFS)

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
		mockFS.EXPECT().Delete(path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(nil)
		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(updatedUser, nil)

		result, err := usecaseAuth.DeleteImage(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, updatedUser, result)
	})

	t.Run("user not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("user not found")

		mockRepo.EXPECT().GetUserByID(ctx, int64(userID)).Return(domain.User{}, expectedErr)

		result, err := usecaseAuth.DeleteImage(ctx, userID)

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
		mockFS.EXPECT().Delete(path.Base(imagePath)).Return(fsErr)

		result, err := usecaseAuth.DeleteImage(ctx, userID)

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
		mockFS.EXPECT().Delete(path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(repoErr)

		result, err := usecaseAuth.DeleteImage(ctx, userID)

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
		mockFS.EXPECT().Delete(path.Base(imagePath)).Return(nil)
		mockRepo.EXPECT().DeleteUserImage(ctx, int64(userID)).Return(nil)
		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(domain.User{}, updateErr)

		result, err := usecaseAuth.DeleteImage(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, updateErr, err)
		assert.Equal(t, domain.User{}, result)
	})
}
