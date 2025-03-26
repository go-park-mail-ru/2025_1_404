package domain

// User Структура пользователя
type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// RegisterRequest Запрос на регистрацию
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest Запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
