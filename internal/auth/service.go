package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Структура для хранения пользователей
var users []User

// IsEmailTaken Проверка, существует ли email
func IsEmailTaken(email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// CreateUser Создание нового пользователя
func CreateUser(email, password, firstName, lastName string) (User, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("ошибка при хешировании пароля")
	}

	user := User{
		ID:        len(users) + 1,
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}
	users = append(users, user)
	return user, nil
}

// GetUserByEmail Получение пользователя по email
func GetUserByEmail(email string) (User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("пользователь не найден")
}

// GetUserByID Получение пользователя по ID
func GetUserByID(id int) (User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("пользователь не найден")
}
