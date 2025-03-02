package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"sync"
)

var users []User         // Список пользователей
var userMutex sync.Mutex // Защита от конкурентного доступа

// Функция хеширования пароля
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Проверка email по регулярному выражению
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Проверка, существует ли пользователь с таким email
func isEmailTaken(email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// Регистрация пользователя с валидацией
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		log.Println("❌ Ошибка: пустое тело запроса")
		sendErrorResponse(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("❌ Ошибка парсинга JSON:", err)
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация email
	if newUser.Email == "" || !isValidEmail(newUser.Email) {
		log.Println("❌ Ошибка регистрации: некорректный email")
		sendErrorResponse(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Проверка "уникальности" email
	userMutex.Lock()
	if isEmailTaken(newUser.Email) {
		userMutex.Unlock()
		log.Println("❌ Ошибка регистрации: email уже используется")
		sendErrorResponse(w, "Email is already taken", http.StatusBadRequest)
		return
	}
	userMutex.Unlock()

	// Валидация пароля (не менее 6 символов)
	if len(newUser.Password) < 6 {
		log.Println("❌ Ошибка регистрации: пароль слишком короткий")
		sendErrorResponse(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Проверяем, что имя и фамилия не пустые
	if newUser.FirstName == "" || newUser.LastName == "" {
		log.Println("❌ Ошибка регистрации: не все поля заполнены")
		sendErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Ошибка хеширования пароля:", err)
		sendErrorResponse(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Добавляем пользователя в список
	userMutex.Lock()
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	userMutex.Unlock()

	log.Printf("✅ Пользователь зарегистрирован: %s (ID: %d)", newUser.Email, newUser.ID)

	sendJSONResponse(w, map[string]interface{}{
		"id":         newUser.ID,
		"email":      newUser.Email,
		"first_name": newUser.FirstName,
		"last_name":  newUser.LastName,
		"is_realtor": newUser.IsRealtor,
	}, http.StatusCreated)
}
