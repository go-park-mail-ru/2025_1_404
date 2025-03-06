package auth

import "errors"

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
func CreateUser(email, password, firstName, lastName string) User {
	user := User{
		ID:        len(users) + 1,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
	users = append(users, user)
	return user
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
