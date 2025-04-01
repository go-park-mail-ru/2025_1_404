package validation

import (
	"testing"
)

// Тест passwordValidator
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Корректный пароль", "Password123", false},
		{"Без заглавных букв", "password123", true},
		{"Без цифр", "Passwordword", true},
	}

	validate := GetValidator()
	type testPassword struct {
		Password string `validate:"password"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := testPassword{Password: tt.password}
			err := validate.Struct(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.password, tt.wantErr, err)
			}
		})
	}
}

// Тест nameValidator
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
		{"Имя с одной буквой", "К.", true},
		{"Пустое имя", "", true},
	}

	type testName struct {
		Name string `validate:"name"`
	}

	validate := GetValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := testName{Name: tt.value}
			err := validate.Struct(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName для %s: ожидалась ошибка = %v, получено ошибка = %v", tt.value, tt.wantErr, err)
			}
		})
	}
}