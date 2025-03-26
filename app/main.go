package main

import (
	delivery "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"log"
	"net/http"
)

// Основной обработчик маршрутов
func main() {
	log.Println("Сервер запущен на http://localhost:8001")
	mux := http.NewServeMux()

	// Обработка всех путей
	mux.HandleFunc("/", utils.NotFoundHandler)

	// Объявления
	mux.HandleFunc("/api/v1/offers", delivery.GetOffersHandler)

	// Авторизация
	mux.HandleFunc("/api/v1/auth/register", delivery.RegisterHandler)
	mux.HandleFunc("/api/v1/auth/login", delivery.LoginHandler)
	mux.HandleFunc("/api/v1/auth/me", delivery.MeHandler)
	mux.HandleFunc("/api/v1/auth/logout", delivery.LogoutHandler)

	// CORS-middleware
	corsMux := middleware.CORSHandler(mux)

	// Запуск сервера
	err := http.ListenAndServe(":8001", corsMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
