package domain

// User Структура пользователя
type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Image     string
}

// RegisterRequest Запрос на регистрацию
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,password,min=8"`
	FirstName string `json:"first_name" validate:"required,name,max=32"`
	LastName  string `json:"last_name" validate:"required,name,max=32"`
}

// LoginRequest Запрос на вход
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UpdateRequest Запрос на обновление данных
type UpdateRequest struct {
	ID        int
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,name,max=32"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,name,max=32"`
}

func UserFromUpdated(updatedUser UpdateRequest) User {
	return User{
		ID:        updatedUser.ID,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
	}
}
