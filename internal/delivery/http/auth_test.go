package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест на регистрацию (успешный)
func TestRegisterHandler_Success(t *testing.T) {
	t.Cleanup(func() { usecase.Users = []domain.User{} })

	reqBody, _ := json.Marshal(domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "Password123",
		FirstName: "Иван",
		LastName:  "Петров",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201, получен %d", rr.Code)
	}
}

// Тест на регистрацию (неправильный метод)
func TestRegisterHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/register", nil)
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус 405, получен %d", rr.Code)
	}
}

// Тест на регистрацию (пустое тело запроса)
func TestRegisterHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{invalid_json}"))
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400, получен %d", rr.Code)
	}
}

// Тест успешного входа
func TestLoginHandler_Success(t *testing.T) {
	t.Cleanup(func() { usecase.Users = []domain.User{} })

	user, _ := usecase.CreateUser("user@example.com", "password", "Иван", "Петров")

	reqBody, _ := json.Marshal(domain.LoginRequest{
		Email:    user.Email,
		Password: "password",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}
}

// Тест на вход (неверный пароль)
func TestLoginHandler_InvalidPassword(t *testing.T) {
	t.Cleanup(func() { usecase.Users = []domain.User{} })

	usecase.CreateUser("user@example.com", "password", "Иван", "Петров") // Обычный пароль

	reqBody, _ := json.Marshal(domain.LoginRequest{
		Email:    "user@example.com",
		Password: "wrongpassword",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rr.Code)
	}
}

// Тест на вход (неверный метод)
func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус 405, получен %d", rr.Code)
	}
}

// Тест получения пользователя (/auth/me с валидным токеном)
func TestMeHandler_Success(t *testing.T) {
	t.Cleanup(func() { usecase.Users = []domain.User{} })

	user, _ := usecase.CreateUser("me@example.com", "securepassword", "Анна", "Сидорова")
	token, _ := utils.GenerateJWT(user.ID)

	req := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	MeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}
}

// Тест /auth/me (нет токена)
func TestMeHandler_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
	rr := httptest.NewRecorder()

	MeHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rr.Code)
	}
}

// Тест /auth/me (неверный токен)
func TestMeHandler_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid_token"})
	rr := httptest.NewRecorder()

	MeHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rr.Code)
	}
}

// Тест логаута (/auth/logout)
func TestLogoutHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rr := httptest.NewRecorder()

	LogoutHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}

	cookie := rr.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Value != "" {
		t.Errorf("Ожидалось удаление токена, но он остался")
	}
}

// Тест /auth/logout (неверный метод)
func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rr := httptest.NewRecorder()

	LogoutHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус 405, получен %d", rr.Code)
	}
}
