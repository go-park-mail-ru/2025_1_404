package main

import (
	"context"
	"log"
	"net/http"

	delivery "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, переменные будут браться из окружения")
	}
}

func main() {
	log.Println("Сервер запущен на http://localhost:8001")

	ctx := context.Background()

	// Инициализация подключения к БД
	dbpool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	// Логгер
	l, _ := logger.NewZapLogger()
	defer l.Close()

	// Репозиторий
	repo := repository.NewRepository(dbpool, l)

	// Юзкейсы
	authUC := usecase.NewAuthUsecase(repo, l)
	offerUC := usecase.NewOfferUsecase(repo, l)

	// Хендлеры
	authHandler := delivery.NewAuthHandler(authUC)
	offerHandler := delivery.NewOfferHandler(offerUC)

	// Маршруты
	mux := http.NewServeMux()

	// Not Found
	mux.HandleFunc("/", utils.NotFoundHandler)

	// Авторизация
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/me", authHandler.Me)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)

	// Объявления
	mux.HandleFunc("/api/v1/offers", offerHandler.GetOffersHandler)

	// AccessLog middleware
	logMux := middleware.AccessLog(l, mux)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux)

	// Запуск сервера
	if err := http.ListenAndServe(":8001", corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
