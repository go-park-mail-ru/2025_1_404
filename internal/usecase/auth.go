package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"golang.org/x/crypto/bcrypt"
)

// Структура для хранения пользователей
var Users []domain.User

// IsEmailTaken Проверка, существует ли email
func IsEmailTaken(email string) bool {
	for _, user := range Users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// CreateUser Создание нового пользователя
func CreateUser(email, password, firstName, lastName string) (domain.User, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, errors.New("ошибка при хешировании пароля")
	}

	user := domain.User{
		ID:        len(Users) + 1,
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}
	Users = append(Users, user)
	return user, nil
}

// GetUserByEmail Получение пользователя по email
func GetUserByEmail(email string) (domain.User, error) {
	for _, user := range Users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, errors.New("пользователь не найден")
}

// GetUserByID Получение пользователя по ID
func GetUserByID(id int) (domain.User, error) {
	for _, user := range Users {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, errors.New("пользователь не найден")
}
