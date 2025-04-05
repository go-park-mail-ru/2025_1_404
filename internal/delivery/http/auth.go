package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/content"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_404/pkg/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UC *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{UC: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	// Валидация email, пароля и имени/фамилии
	validate := validation.GetValidator()

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

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	// Проверка полей запроса
	validate := validation.GetValidator()

	err := validate.Struct(req)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest)
		return
	}

	// Ищем юзера по почте
	user, err := h.UC.GetUserByEmail(r.Context(), req.Email)

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

func (h *AuthHandler) Update(w http.ResponseWriter, r *http.Request) {

	// Достаем из контекста от Auth middleware id юзера
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	var updateUser domain.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	updateUser.ID = userID

	// Валидируем те поля, которые передали
	validate := validation.GetValidator()

	err := validate.Struct(updateUser)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest)
		return
	}

	// Возвращает User вместо UserUpdate
	user := domain.UserFromUpdated(updateUser)
	// Пытаемся обновить данные
	userUpdated, err := h.UC.UpdateUser(r.Context(), user)
	if err != nil {
		utils.SendErrorResponse(w, "не удалось обновить данные о пользователе", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, userUpdated, http.StatusOK)
}

func (h *AuthHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	// Достаем из контекста от Auth middleware id юзера
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "Upload image failed", http.StatusBadRequest)
		return
	}

	// Достаем file по name="avatar"
	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.SendErrorResponse(w, "Failed to get image from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(w, "invalid content type", http.StatusBadRequest)
		return
	}

	// Проверяем подходит ли файл под условия и получаем его расширение
	contentType, err := content.CheckImage(fileBytes)
	if err != nil {
		utils.SendErrorResponse(w, "invalid file type or size", http.StatusBadRequest)
		return
	}

	upload := filestorage.FileUpload{
		Name:        uuid.New().String()+"."+contentType,
		Size:        header.Size,
		File:        bytes.NewReader(fileBytes),
		ContentType: contentType,
	}

	// Пытаемся загрузить картинку
	updatedUser, err := h.UC.UploadImage(r.Context(), userID, upload)
	if err != nil {
		utils.SendErrorResponse(w, "failed to upload image", http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, updatedUser, http.StatusOK)
}
