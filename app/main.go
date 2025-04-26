package main

import (
	"context"
	"log"
	"net/http"

	delivery "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http"
	repository "github.com/go-park-mail-ru/2025_1_404/internal/repository/csat"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, переменные будут браться из окружения")
	}
}

func main() {
	log.Println("Сервер запущен на ", utils.BasePath)

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

	csatRepo := repository.NewCsatRepository(dbpool, l)
	csatUC := usecase.NewCsatUsecase(csatRepo, l)
	csatHandler := delivery.NewCsatHandler(csatUC)

	// Маршруты
	r := mux.NewRouter()

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	r.HandleFunc("/api/v1/csat", csatHandler.GetQuestionsByEvent).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/csat", csatHandler.AddAnswerToQuestion).
		Methods(http.MethodPost)
	r.Handle("/api/v1/csat/stats", middleware.AuthHandler(l, http.HandlerFunc(csatHandler.GetAnswersByQuestion))).
		Methods(http.MethodGet)

	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux)

	// Запуск сервера
	if err := http.ListenAndServe(":8002", corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
