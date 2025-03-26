package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_404/pkg/validation"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// RegisterHandler Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	// Проверка полей запроса
	if err := validation.ValidateRegisterRequest(req); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация email, пароля и имени/фамилии
	if err := validation.ValidateUser(domain.User{
		0,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	}); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка уникальности email
	if usecase.IsEmailTaken(req.Email) {
		utils.SendErrorResponse(w, "Email уже занят", http.StatusBadRequest)
		return
	}

	// Создаём пользователя
	user, err := usecase.CreateUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SendJSONResponse(w, map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}, http.StatusCreated)
}

// LoginHandler Авторизация пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	// Проверка полей запроса
	if err := validation.ValidateLoginRequest(req); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем юзера по почте
	user, err := usecase.GetUserByEmail(req.Email)
	if err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	// Генерируем JWT
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Устанавливаем JWT в куки
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SendJSONResponse(w, map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}, http.StatusOK)
}

// MeHandler Получение текущего пользователя
func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		utils.SendErrorResponse(w, "Учётные данные не предоставлены", http.StatusUnauthorized)
		return
	}

	claims, err := utils.ParseJWT(cookie.Value)
	if err != nil {
		utils.SendErrorResponse(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	user, err := usecase.GetUserByID(claims.UserID)
	if err != nil {
		utils.SendErrorResponse(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	utils.SendJSONResponse(w, map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}, http.StatusOK)
}

// LogoutHandler Логаут пользователя (удаление куки с JWT)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "", // Пустой токен
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Кука сразу "просроченная"
	})

	utils.SendJSONResponse(w, map[string]string{"message": "Успешный выход"}, http.StatusOK)
}
