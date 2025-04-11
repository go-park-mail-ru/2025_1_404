package main

import (
	"context"
	"log"
	"net/http"

	delivery "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database"
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

	// Хранилище файлов
	basePath := "./../internal/static/upload"
	fs := filestorage.NewLocalStorage(basePath)

	// Репозиторий
	repo := repository.NewRepository(dbpool, l)

	// Юзкейсы
	authUC := usecase.NewAuthUsecase(repo, l, fs)
	offerUC := usecase.NewOfferUsecase(repo, l, fs)

	// Хендлеры
	authHandler := delivery.NewAuthHandler(authUC)
	offerHandler := delivery.NewOfferHandler(offerUC)

	// Маршруты
	r := mux.NewRouter()

	// Static
	r.PathPrefix("/images/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filestorage.ServeFile(w, r, basePath)
	}))

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	// Авторизация
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).
		Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).
		Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/logout", authHandler.Logout).
		Methods(http.MethodPost)

	// Профиль
	r.Handle("/api/v1/auth/me", middleware.AuthHandler(l, http.HandlerFunc(authHandler.Me))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/users/update", middleware.AuthHandler(l, http.HandlerFunc(authHandler.Update))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image", middleware.AuthHandler(l, http.HandlerFunc(authHandler.UploadImage))).
		Methods(http.MethodPut)

	// Объявления
	r.HandleFunc("/api/v1/offers", offerHandler.GetOffersHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/offers/{id:[0-9]+}", offerHandler.GetOfferByID).
		Methods(http.MethodGet)
	r.Handle("/api/v1/offers", middleware.AuthHandler(l, http.HandlerFunc(offerHandler.CreateOffer))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}", middleware.AuthHandler(l, http.HandlerFunc(offerHandler.UpdateOffer))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/offers/{id:[0-9]+}", middleware.AuthHandler(l, http.HandlerFunc(offerHandler.DeleteOffer))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/offers/{id:[0-9]+}/publish", middleware.AuthHandler(l, http.HandlerFunc(offerHandler.PublishOffer))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}/image", middleware.AuthHandler(l, http.HandlerFunc(offerHandler.UploadOfferImage))).
		Methods(http.MethodPost)

	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux)

	// Запуск сервера
	if err := http.ListenAndServe(":8001", corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
