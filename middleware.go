package main

import (
	"log"
	"net/http"
	"time"
)

// ResponseWriter с логированием статуса
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Реализация WriteHeader с логированием кода ответа
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Middleware для логирования всех HTTP-запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter для логирования статуса
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		// Логируем запрос: IP, метод, путь, статус и время обработки
		log.Printf("%s %s %s %d %v", r.RemoteAddr, r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}
