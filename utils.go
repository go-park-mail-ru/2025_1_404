package main

import (
	"encoding/json"
	"net/http"
)

// Отправка успешного JSON-ответа
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Отправка ошибки в JSON-формате
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// CORS Middleware (чтобы фронтенд мог обращаться к API)
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func corsMiddleware(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}
