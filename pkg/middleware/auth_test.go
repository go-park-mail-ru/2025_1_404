package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func userIDHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot) 
	id := r.Context().Value(utils.UserIDKey).(int)
	w.Header().Set("id", strconv.Itoa(id))
}

func TestAuthMiddleware_OK(t *testing.T) {
	cookie, _ := utils.GenerateJWT(1)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Cookie", fmt.Sprintf(`token=%s`, cookie))
	rr := httptest.NewRecorder()
	l, _ := logger.NewZapLogger()

	handler := AuthHandler(l, http.HandlerFunc(userIDHandler))
	handler.ServeHTTP(rr, req)	

	expectedID := "1"
	actualID := rr.Header().Get("id")

	if rr.Code != http.StatusTeapot {
		t.Errorf("Ожидался статус 418, получен %d", rr.Code)
	}

	if expectedID != actualID {
		t.Errorf("Ожидался id %s, получен: %s", expectedID, actualID)
	}
}

func TestAuthMiddleware_Fail_EmptyCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	l, _ := logger.NewZapLogger()

	handler := AuthHandler(l, http.HandlerFunc(userIDHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rr.Code)
	}

	message := "Учётные данные не предоставлены"

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if (response["error"]) != message {
		t.Errorf(`Ожидалось сообщение "%s", получено ""%s`, message, response["error"])
	}
}

func TestAuthMiddleware_Fail_IncorrectToken(t *testing.T) {
	cookie := "badCookie"

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Cookie", fmt.Sprintf(`token=%s`, cookie))
	rr := httptest.NewRecorder()
	l, _ := logger.NewZapLogger()

	handler := AuthHandler(l, http.HandlerFunc(userIDHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rr.Code)
	}

	message := "Неверный токен"

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if (response["error"]) != message {
		t.Errorf(`Ожидалось сообщение "%s", получено ""%s`, message, response["error"])
	}
}

