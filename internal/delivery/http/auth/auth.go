package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/content"
	"github.com/go-park-mail-ru/2025_1_404/pkg/csrf"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_404/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UC  authUsecase
	cfg *config.Config
}

func NewAuthHandler(uc authUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{UC: uc, cfg: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Валидация email, пароля и имени/фамилии
	validate := validation.GetValidator()

	err := validate.Struct(req)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	user, err := h.UC.CreateUser(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании токена", http.StatusInternalServerError, &h.cfg.App.CORS)
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

	utils.SendJSONResponse(w, user, http.StatusCreated, &h.cfg.App.CORS)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Проверка полей запроса
	validate := validation.GetValidator()

	err := validate.Struct(req)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Ищем юзера по почте
	user, err := h.UC.GetUserByEmail(r.Context(), req.Email)

	if err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized, &h.cfg.App.CORS)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendErrorResponse(w, "Неверная почта или пароль", http.StatusUnauthorized, &h.cfg.App.CORS)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании токена", http.StatusInternalServerError, &h.cfg.App.CORS)
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

	utils.SendJSONResponse(w, user, http.StatusOK, &h.cfg.App.CORS)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	user, err := h.UC.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.SendErrorResponse(w, "Пользователь не найден", http.StatusUnauthorized, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, user, http.StatusOK, &h.cfg.App.CORS)
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

	utils.SendJSONResponse(w, map[string]string{"message": "Успешный выход"}, http.StatusOK, &h.cfg.App.CORS)
}

func (h *AuthHandler) Update(w http.ResponseWriter, r *http.Request) {

	// Достаем из контекста от Auth middleware id юзера
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	var updateUser domain.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	updateUser.ID = userID

	// Валидируем те поля, которые передали
	validate := validation.GetValidator()

	err := validate.Struct(updateUser)
	if err != nil {
		utils.SendErrorResponse(w, validation.GetError(err), http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Возвращает User вместо UserUpdate
	user := domain.UserFromUpdated(updateUser)
	// Пытаемся обновить данные
	userUpdated, err := h.UC.UpdateUser(r.Context(), user)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, userUpdated, http.StatusOK, &h.cfg.App.CORS)
}

func (h *AuthHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	// Достаем из контекста от Auth middleware id юзера
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "Upload image failed", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Достаем file по name="avatar"
	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.SendErrorResponse(w, "Failed to get image from request", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(w, "invalid content type", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	// Проверяем подходит ли файл под условия и получаем его расширение
	contentType, err := content.CheckImage(fileBytes)
	if err != nil {
		utils.SendErrorResponse(w, "invalid file type or size", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	upload := s3.Upload{
		Bucket:      "avatars",
		Filename:    header.Filename,
		Size:        header.Size,
		File:        bytes.NewReader(fileBytes),
		ContentType: contentType,
	}

	// Пытаемся загрузить картинку
	updatedUser, err := h.UC.UploadImage(r.Context(), userID, upload)
	if err != nil {
		utils.SendErrorResponse(w, "failed to upload image", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, updatedUser, http.StatusOK, &h.cfg.App.CORS)
}

func (h *AuthHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "Не удалось удалить фотографию", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	updatedUser, err := h.UC.DeleteImage(r.Context(), userID)
	if err != nil {
		utils.SendErrorResponse(w, "Не удалось удалить фотографию", http.StatusInternalServerError, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, updatedUser, http.StatusOK, &h.cfg.App.CORS)
}

func (h *AuthHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int)

	token := csrf.GenerateCSRF(strconv.Itoa(userID), h.cfg.App.Auth.CSRF.Salt)

	response := csrf.GetCSRFResponse(token)
	utils.SendJSONResponse(w, response, http.StatusOK, &h.cfg.App.CORS)
}
