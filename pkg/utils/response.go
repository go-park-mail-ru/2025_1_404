package utils

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse Отправка успешного JSON-ответа
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendErrorResponse Отправка ошибки в JSON-формате
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// EnableCORS CORS Middleware (чтобы фронтенд мог обращаться к API)
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", BaseFrontendPath)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-csrf-token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// NotFoundHandler 404 на несуществующие урлы
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
