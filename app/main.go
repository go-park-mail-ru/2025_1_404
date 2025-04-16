package main

import (
	"context"
	"log"
	"net/http"

	deliveryAuth "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http/auth"
	deliveryOffer "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http/offer"
	deliveryZhk "github.com/go-park-mail-ru/2025_1_404/internal/delivery/http/zhk"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	repoAuth "github.com/go-park-mail-ru/2025_1_404/internal/repository/auth"
	repoOffer "github.com/go-park-mail-ru/2025_1_404/internal/repository/offer"
	repoZhk "github.com/go-park-mail-ru/2025_1_404/internal/repository/zhk"
	usecaseAuth "github.com/go-park-mail-ru/2025_1_404/internal/usecase/auth"
	usecaseOffer "github.com/go-park-mail-ru/2025_1_404/internal/usecase/offer"
	usecaseZhk "github.com/go-park-mail-ru/2025_1_404/internal/usecase/zhk"
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

	// Хранилище файлов
	basePath := "./../internal/static/upload"
	fs := filestorage.NewLocalStorage(basePath)

	// Репозиторий
	authRepo := repoAuth.NewAuthRepository(dbpool, l)
	offerRepo := repoOffer.NewOfferRepository(dbpool, l)
	zhkRepo := repoZhk.NewZhkRepository(dbpool, l)

	// Юзкейсы
	authUC := usecaseAuth.NewAuthUsecase(authRepo, l, fs)
	offerUC := usecaseOffer.NewOfferUsecase(offerRepo, l, fs)
	zhkUC := usecaseZhk.NewZhkUsecase(zhkRepo, l)

	// Хендлеры
	authHandler := deliveryAuth.NewAuthHandler(authUC)
	offerHandler := deliveryOffer.NewOfferHandler(offerUC)
	zhkHandler := deliveryZhk.NewZhkHandler(zhkUC)

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
	r.Handle("/api/v1/users/update",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(authHandler.Update)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(authHandler.UploadImage)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(authHandler.DeleteImage)))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/users/csrf", middleware.AuthHandler(l, http.HandlerFunc(authHandler.GetCSRFToken))).
		Methods(http.MethodGet)

	// Объявления
	r.HandleFunc("/api/v1/offers", offerHandler.GetOffersHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/offers/{id:[0-9]+}", offerHandler.GetOfferByID).
		Methods(http.MethodGet)
	r.Handle("/api/v1/offers",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.CreateOffer)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.UpdateOffer)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/offers/{id:[0-9]+}",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.DeleteOffer)))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/offers/{id:[0-9]+}/publish",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.PublishOffer)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}/image",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.UploadOfferImage)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/images/{id:[0-9]+}",
		middleware.AuthHandler(l, middleware.CSRFMiddleware(l, http.HandlerFunc(offerHandler.DeleteOfferImage)))).
		Methods(http.MethodDelete)

	// ЖК
	r.HandleFunc("/api/v1/zhk/{id:[0-9]+}", zhkHandler.GetZhkInfo).Methods("GET")
	r.HandleFunc("/api/v1/zhks", zhkHandler.GetAllZhk).Methods("GET")

	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux)

	// Запуск сервера
	if err := http.ListenAndServe(":8001", corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
