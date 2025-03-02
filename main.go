package main

import (
	"log"
	"net/http"
)

// Основной обработчик маршрутов
func main() {
	log.Println("Сервер запущен на http://localhost:8000")
	mux := http.NewServeMux()

	// Обработка OPTIONS-запросов
	mux.HandleFunc("/", corsMiddleware)
	// Эндпоинт для объявлений
	mux.HandleFunc("/offers", getOffers)
	// Эндпоинт для регистрации пользователей
	mux.HandleFunc("/auth/register/", registerUser)

	// Добавляем логирование всех запросов
	loggedMux := loggingMiddleware(mux)

	// Запуск сервера
	err := http.ListenAndServe(":8000", loggedMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
