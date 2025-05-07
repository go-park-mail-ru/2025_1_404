package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
)

// Тест SendJSONResponse (корректный JSON-ответ)
func TestSendJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "OK"}
	cfg := &config.CORSConfig{AllowOrigin: "http://localhost:8000", AllowMethods: "GET, POST, PUT, OPTIONS, DELETE",
		AllowHeaders: "Content-Type, x-csrf-token", AllowCredentials: "true"}

	SendJSONResponse(rr, data, http.StatusOK, cfg)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получено %d", rr.Code)
	}

	// Проверяем JSON-ответ
	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка сериализации в JSON: %v", err)
	}

	if response["message"] != "OK" {
		t.Errorf("Ожидалось сообщение 'OK', получено '%s'", response["message"])
	}

	// Проверяем заголовки
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Ожидался Content-Type 'application/json', got '%s'", rr.Header().Get("Content-Type"))
	}
}

// Тест SendErrorResponse (ошибочный JSON-ответ)
func TestSendErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	errorMessage := "Something went wrong"
	cfg := &config.CORSConfig{AllowOrigin: "http://localhost:8000", AllowMethods: "GET, POST, PUT, OPTIONS, DELETE",
		AllowHeaders: "Content-Type, x-csrf-token", AllowCredentials: "true"}

	SendErrorResponse(rr, errorMessage, http.StatusBadRequest, cfg)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400, получено %d", rr.Code)
	}

	// Проверяем JSON-ответ
	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка сериализации в JSON: %v", err)
	}

	if response["error"] != errorMessage {
		t.Errorf("Ожидалось сообщение '%s', получено '%s'", errorMessage, response["error"])
	}

	// Проверяем заголовки
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Ожидался Content-Type 'application/json', получено '%s'", rr.Header().Get("Content-Type"))
	}
}

// Тест EnableCORS (CORS-заголовки)
func TestEnableCORS(t *testing.T) {
	rr := httptest.NewRecorder()

	cfg := &config.CORSConfig{
		AllowOrigin:      "http://localhost:8000",
		AllowMethods:     "GET, POST, PUT, OPTIONS, DELETE",
		AllowHeaders:     "Content-Type, x-csrf-token",
		AllowCredentials: "true",
	}

	EnableCORS(rr, cfg)

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "http://localhost:8000",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, x-csrf-token",
		"Access-Control-Allow-Credentials": "true",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("Ожидался хедер %s: '%s', получено '%s'", header, expectedValue, value)
		}
	}
}

// Тест NotFoundHandler (404-ответ)
func TestNotFoundHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/no-route", nil)
	rr := httptest.NewRecorder()

	NotFoundHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404, получено %d", rr.Code)
	}

	expectedBody := "404 Not Found\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Ожидалось сообщение '%s', получено '%s'", expectedBody, rr.Body.String())
	}
}
