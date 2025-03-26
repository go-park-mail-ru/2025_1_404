package validation

import (
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"testing"
)

// Тест ValidateEmail
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Корректный email", "user@example.com", false},
		{"Корректный email с чисоами", "test123@mail.ru", false},
		{"Некорректный email 1", "bademail", true},
		{"Некорректный email 2", "user@com", true},
		{"Пустой email", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.email, tt.wantErr, err)
			}
		})
	}
}

// Тест ValidatePassword
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Корректный пароль", "Password123", false},
		{"Короткий пароль", "Pass1", true},
		{"Без заглавных букв", "password123", true},
		{"Без цифр", "Passwordword", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.password, tt.wantErr, err)
			}
		})
	}
}

// Тест ValidateName
func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"Корректное имя", "Вадим", false},
		{"Корректная фамилия", "Иванов", false},
		{"Имя с тире", "Петров-Сидоров", false},
		{"Имя с цифрами", "вадим123", true},
		{"Слишком длинное имя", "Изподвыподвывертовичкин", true},
		{"Имя с одной буквой", "К.", true},
		{"Пустое имя", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.value, tt.wantErr, err)
			}
		})
	}
}

// Тест ValidateRegisterRequest (поля регистрации)
func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name    string
		request domain.RegisterRequest
		wantErr bool
	}{
		{
			name: "Валидный запрос",
			request: domain.RegisterRequest{
				Email:     "user@example.com",
				Password:  "SecurePass123",
				FirstName: "Иван",
				LastName:  "Петров",
			},
			wantErr: false,
		},
		{
			name: "Пустой email",
			request: domain.RegisterRequest{
				Email:     "",
				Password:  "SecurePass123",
				FirstName: "Иван",
				LastName:  "Петров",
			},
			wantErr: true,
		},
		{
			name: "Пустой пароль",
			request: domain.RegisterRequest{
				Email:     "user@example.com",
				Password:  "",
				FirstName: "Иван",
				LastName:  "Петров",
			},
			wantErr: true,
		},
		{
			name: "Пустое имя",
			request: domain.RegisterRequest{
				Email:     "user@example.com",
				Password:  "SecurePass123",
				FirstName: "",
				LastName:  "Петров",
			},
			wantErr: true,
		},
		{
			name: "Пустая фамилия",
			request: domain.RegisterRequest{
				Email:     "user@example.com",
				Password:  "SecurePass123",
				FirstName: "Иван",
				LastName:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegisterRequest для '%s': ожидалась ошибка = %v, получено ошибка = %v", tt.name, tt.wantErr, err)
			}
		})
	}
}

// Тест ValidateLoginRequest (поля логина)
func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name    string
		request domain.LoginRequest
		wantErr bool
	}{
		{
			name: "Корректный запрос",
			request: domain.LoginRequest{
				Email:    "user@example.com",
				Password: "Password123",
			},
			wantErr: false,
		},
		{
			name: "Пустой email",
			request: domain.LoginRequest{
				Email:    "",
				Password: "Password123",
			},
			wantErr: true,
		},
		{
			name: "Пустой пароль",
			request: domain.LoginRequest{
				Email:    "user@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoginRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoginRequest для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.name, tt.wantErr, err)
			}
		})
	}
}
