package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
)

// SendJSONResponse Отправка успешного JSON-ответа
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int, cfg *config.CORSConfig) {
	EnableCORS(w, cfg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendErrorResponse Отправка ошибки в JSON-формате
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int, cfg *config.CORSConfig) {
	EnableCORS(w, cfg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// EnableCORS CORS Middleware (чтобы фронтенд мог обращаться к API)
func EnableCORS(w http.ResponseWriter, cfg *config.CORSConfig) {
	w.Header().Set("Access-Control-Allow-Origin", cfg.AllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
	w.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
	w.Header().Set("Access-Control-Allow-Credentials", cfg.AllowCredentials)
}

// NotFoundHandler 404 на несуществующие урлы
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
