package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_404/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UC *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{UC: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	// Валидация email, пароля и имени/фамилии
	validate := validation.GetValidator();

	err := validate.Struct(req)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest)
		return
	}

	if h.UC.IsEmailTaken(r.Context(), req.Email) {
		utils.SendErrorResponse(w, "Email уже занят", http.StatusBadRequest)
		return
	}

	user, err := h.UC.CreateUser(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
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

	utils.SendJSONResponse(w, user, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	validate := validation.GetValidator();

	err := validate.Struct(req)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest)
		return
	}
	
	// Ищем юзера по почте
	user, err := usecase.GetUserByEmail(req.Email)

	if err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
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

	utils.SendJSONResponse(w, user, http.StatusOK)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	user, err := h.UC.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.SendErrorResponse(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	utils.SendJSONResponse(w, user, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})

	utils.SendJSONResponse(w, map[string]string{"message": "Успешный выход"}, http.StatusOK)
}
