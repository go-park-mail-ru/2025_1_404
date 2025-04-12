package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	mockFS "github.com/go-park-mail-ru/2025_1_404/internal/filestorage/mocks"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	mockRepo "github.com/go-park-mail-ru/2025_1_404/internal/repository/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	//"golang.org/x/crypto/bcrypt"
)

// // Тест IsEmailTaken (существующий и новый email)
// func TestIsEmailTaken(t *testing.T) {
// 	t.Cleanup(func() { Users = []domain.User{} })

// 	CreateUser("user@example.com", "password", "Иван", "Петров")

// 	if !IsEmailTaken("user@example.com") {
// 		t.Errorf("Ожидалось, что email занят, но функция вернула false")
// 	}

// 	if IsEmailTaken("new@example.com") {
// 		t.Errorf("Ожидалось, что email свободен, но функция вернула true")
// 	}
// }

// // Тест CreateUser (успешное создание пользователя)
// func TestCreateUser_Success(t *testing.T) {
// 	t.Cleanup(func() { Users = []domain.User{} })

// 	user, _ := CreateUser("test@example.com", "SecurePass123", "Анна", "Смирнова")

// 	if user.ID != 1 {
// 		t.Errorf("Ожидался ID = 1, получен %d", user.ID)
// 	}
// 	if user.Email != "test@example.com" {
// 		t.Errorf("Ожидался Email = test@example.com, получен %s", user.Email)
// 	}
// 	if user.FirstName != "Анна" {
// 		t.Errorf("Ожидалось FirstName = Анна, получено %s", user.FirstName)
// 	}
// 	if user.LastName != "Смирнова" {
// 		t.Errorf("Ожидалось LastName = Смирнова, получено %s", user.LastName)
// 	}

// 	// Проверяем, что пароль хеширован
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("SecurePass123")); err != nil {
// 		t.Errorf("Пароль не прошел проверку хеша")
// 	}
// }

// // Тест GetUserByEmail (существующий и несуществующий email)
// func TestGetUserByEmail(t *testing.T) {
// 	t.Cleanup(func() { Users = []domain.User{} })

// 	CreateUser("user@example.com", "password", "Иван", "Петров")

// 	// Поиск существующего пользователя
// 	user, err := GetUserByEmail("user@example.com")
// 	if err != nil {
// 		t.Fatalf("Ошибка при поиске пользователя: %v", err)
// 	}
// 	if user.Email != "user@example.com" {
// 		t.Errorf("Ожидался email user@example.com, получен %s", user.Email)
// 	}

// 	// Поиск несуществующего пользователя
// 	_, err = GetUserByEmail("notfound@example.com")
// 	if err == nil {
// 		t.Errorf("Ожидалась ошибка 'пользователь не найден'")
// 	}
// }

// // Тест GetUserByID (существующий и несуществующий ID)
// func TestGetUserByID(t *testing.T) {
// 	t.Cleanup(func() { Users = []domain.User{} })

// 	user, _ := CreateUser("id@example.com", "password", "ID", "Тест")

// 	// Поиск существующего пользователя
// 	foundUser, err := GetUserByID(user.ID)
// 	if err != nil {
// 		t.Fatalf("Ошибка при поиске пользователя: %v", err)
// 	}
// 	if foundUser.ID != user.ID {
// 		t.Errorf("Ожидался ID %d, получен %d", user.ID, foundUser.ID)
// 	}

// 	// Поиск несуществующего пользователя
// 	_, err = GetUserByID(999)
// 	if err == nil {
// 		t.Errorf("Ожидалась ошибка 'пользователь не найден'")
// 	}
// }

func TestUploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
    defer ctrl.Finish()
	
	mockRepo := mockRepo.NewMockRepository(ctrl)
	mockFS := mockFS.NewMockFileStorage(ctrl)
	logger, _ := logger.NewZapLogger()

	usecase := NewAuthUsecase(mockRepo, logger, mockFS)

	file := filestorage.FileUpload {
		Name: "avatar.png",
		Size: 1024,
		File: nil,
		ContentType: "png",
	}

	user := domain.User {ID: 1, Image: "avatar.png"}

	mockRepo.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(repository.User{ID: 1}, nil).Times(2)
	mockFS.EXPECT().Add(file).Return(nil)
	mockRepo.EXPECT().CreateImage(gomock.Any(), file).Return(nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any(), domain.User{ID: 1, Image: file.Name}).Return(user, nil)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "333")
	updatedUser, err := usecase.UploadImage(ctx, 1, file)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, updatedUser.ID)
	assert.Equal(t, user.Image, updatedUser.Image)
}