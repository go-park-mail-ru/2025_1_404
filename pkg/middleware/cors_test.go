package middleware

import (
	"github.com/go-park-mail-ru/2025_1_404/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Фейковый обработчик для тестов
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot) // 418 I'm a teapot :):)
}

// Тест CORSHandler (OPTIONS-запрос)
func TestCORSHandler_OptionsRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rr := httptest.NewRecorder()

	cfg := &config.CORSConfig{
		AllowOrigin:      "http://localhost:8000",
		AllowMethods:     "GET, POST, PUT, OPTIONS, DELETE",
		AllowHeaders:     "Content-Type, x-csrf-token",
		AllowCredentials: "true",
	}
	handler := CORSHandler(http.HandlerFunc(dummyHandler), cfg)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}

	// Проверяем, что заголовки CORS установлены
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "http://localhost:8000",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, x-csrf-token",
		"Access-Control-Allow-Credentials": "true",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("Ожидался заголовок %s = '%s', получен '%s'", header, expectedValue, value)
		}
	}
}

// Тест CORSHandler (GET-запрос)
func TestCORSHandler_GetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	cfg := &config.CORSConfig{
		AllowOrigin:      "http://localhost:8000",
		AllowMethods:     "GET, POST, PUT, OPTIONS, DELETE",
		AllowHeaders:     "Content-Type, x-csrf-token",
		AllowCredentials: "true",
	}
	handler := CORSHandler(http.HandlerFunc(dummyHandler), cfg)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTeapot { // Должен пройти в dummyHandler
		t.Errorf("Ожидался статус 418, получен %d", rr.Code)
	}

	// Проверяем, что заголовки CORS установлены
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "http://localhost:8000",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, x-csrf-token",
		"Access-Control-Allow-Credentials": "true",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("Ожидался заголовок %s = '%s', получен '%s'", header, expectedValue, value)
		}
	}
}
