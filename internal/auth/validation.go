package auth

import (
	"errors"
	"regexp"
	"unicode"
)

// Регулярка для email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameRegex = regexp.MustCompile(`^[A-Za-zА-Яа-яЁё-]+$`)

// ValidateEmail Валидация email
func ValidateEmail(email string) error {
	if email == "" || !emailRegex.MatchString(email) {
		return errors.New("неверный формат e-mail")
	}
	return nil
}

// ValidatePassword Валидация пароля
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("пароль должен быть длиннее 8 символов")
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("пароль должен включать хотя бы одну букву каждого регистра и цифру")
	}

	return nil
}

func ValidateName(name string) error {
	if name == "" || !nameRegex.MatchString(name) || len(name) > 32 {
		return errors.New("неверное имя/фамилия")
	}

	return nil
}

// ValidateRegisterRequest Проверка полей регистрации
func ValidateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return errors.New("все поля обязательны")
	}
	return nil
}

// ValidateLoginRequest Проверка полей логина
func ValidateLoginRequest(req LoginRequest) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("все поля обязательны")
	}
	return nil
}
