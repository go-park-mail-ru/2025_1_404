package main

import (
	"github.com/go-park-mail-ru/2025_1_404/internal/auth"
	"github.com/go-park-mail-ru/2025_1_404/internal/offers"
	"github.com/go-park-mail-ru/2025_1_404/middleware"
	"github.com/go-park-mail-ru/2025_1_404/utils"
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
	mux.HandleFunc("/api/v1/offers", offers.GetOffersHandler)

	// Авторизация
	mux.HandleFunc("/api/v1/auth/register", auth.RegisterHandler)
	mux.HandleFunc("/api/v1/auth/login", auth.LoginHandler)
	mux.HandleFunc("/api/v1/auth/me", auth.MeHandler)
	mux.HandleFunc("/api/v1/auth/logout", auth.LogoutHandler)

	// CORS-middleware
	corsMux := middleware.CORSHandler(mux)

	// Запуск сервера
	err := http.ListenAndServe(":8001", corsMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
