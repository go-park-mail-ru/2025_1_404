package auth

import "testing"

// Тест валидации почты
func TestValidateEmail(t *testing.T) {
	validEmails := []string{"user@example.com", "test123@mail.ru"}
	invalidEmails := []string{"bademail", "user@com", ""}

	for _, email := range validEmails {
		if err := ValidateEmail(email); err != nil {
			t.Errorf("Верная почта не прошла валидацию: %s", email)
		}
	}

	for _, email := range invalidEmails {
		if err := ValidateEmail(email); err == nil {
			t.Errorf("Неверная почта прошла валидацию: %s", email)
		}
	}
}

// Тест валидации пароля
func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("1234567"); err == nil {
		t.Errorf("Ожидалась ошибка для короткого пароля")
	}
	if err := ValidatePassword("password123"); err == nil {
		t.Errorf("Ожидалась ошибка об отсутствии верхнего регистра")
	}
}

// Тест ValidateRegisterRequest (поля регистрации)
func TestValidateRegisterRequest(t *testing.T) {
	validReq := RegisterRequest{
		Email:     "user@example.com",
		Password:  "SecurePass123",
		FirstName: "Иван",
		LastName:  "Петров",
	}

	invalidReqs := []RegisterRequest{
		{"", "SecurePass123", "Иван", "Петров"},             // Пустой email
		{"user@example.com", "", "Иван", "Петров"},          // Пустой пароль
		{"user@example.com", "SecurePass123", "", "Петров"}, // Пустое имя
		{"user@example.com", "SecurePass123", "Иван", ""},   // Пустая фамилия
	}

	// Валидный запрос должен пройти
	if err := ValidateRegisterRequest(validReq); err != nil {
		t.Errorf("Ожидалось, что запрос с корректными данными пройдет, но получена ошибка: %v", err)
	}

	// Некорректные запросы должны вызывать ошибку
	for _, req := range invalidReqs {
		if err := ValidateRegisterRequest(req); err == nil {
			t.Errorf("Ожидалась ошибка для некорректного запроса %+v, но её не было", req)
		}
	}
}

// Тест ValidateLoginRequest (поля логина)
func TestValidateLoginRequest(t *testing.T) {
	validReq := LoginRequest{
		Email:    "user@example.com",
		Password: "SecurePass123",
	}

	invalidReqs := []LoginRequest{
		{"", "SecurePass123"},    // Пустой email
		{"user@example.com", ""}, // Пустой пароль
	}

	// Валидный запрос должен пройти
	if err := ValidateLoginRequest(validReq); err != nil {
		t.Errorf("Ожидалось, что запрос с корректными данными пройдет, но получена ошибка: %v", err)
	}

	// Некорректные запросы должны вызывать ошибку
	for _, req := range invalidReqs {
		if err := ValidateLoginRequest(req); err == nil {
			t.Errorf("Ожидалась ошибка для некорректного запроса %+v, но её не было", req)
		}
	}
}
