package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
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
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("не удалось загрузить конфиг: %v", err)
	}

	log.Println("Сервер запущен на ", cfg.App.BaseDir)

	ctx := context.Background()

	// Инициализация подключения к БД
	dbpool, err := database.NewPool(&cfg.Postgres, ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	// Логгер
	l, _ := logger.NewZapLogger()
	defer l.Close()

	// Хранилище файлов
	basePath := "./../internal/static/upload"

	s3repo, err := s3.New(&cfg.Minio, l)
	if err != nil {
		log.Printf("не удалось подключиться к s3: %v", err)
		return
	}

	// Репозиторий
	authRepo := repoAuth.NewAuthRepository(dbpool, l)
	offerRepo := repoOffer.NewOfferRepository(dbpool, l)
	zhkRepo := repoZhk.NewZhkRepository(dbpool, l)

	// Юзкейсы
	authUC := usecaseAuth.NewAuthUsecase(authRepo, l, s3repo, cfg)
	offerUC := usecaseOffer.NewOfferUsecase(offerRepo, l, s3repo, cfg)
	zhkUC := usecaseZhk.NewZhkUsecase(zhkRepo, l, cfg)

	// Хендлеры
	authHandler := deliveryAuth.NewAuthHandler(authUC, cfg)
	offerHandler := deliveryOffer.NewOfferHandler(offerUC, cfg)
	zhkHandler := deliveryZhk.NewZhkHandler(zhkUC, cfg)

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
	r.Handle("/api/v1/auth/me", middleware.AuthHandler(l, &cfg.App.CORS,http.HandlerFunc(authHandler.Me))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/users/update",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l,  cfg, http.HandlerFunc(authHandler.Update)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(authHandler.UploadImage)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(authHandler.DeleteImage)))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/users/csrf", middleware.AuthHandler(l, &cfg.App.CORS, http.HandlerFunc(authHandler.GetCSRFToken))).
		Methods(http.MethodGet)

	// Объявления
	r.HandleFunc("/api/v1/offers", offerHandler.GetOffersHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/offers/{id:[0-9]+}", offerHandler.GetOfferByID).
		Methods(http.MethodGet)
	r.Handle("/api/v1/offers",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.CreateOffer)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.UpdateOffer)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/offers/{id:[0-9]+}",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.DeleteOffer)))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/offers/{id:[0-9]+}/publish",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.PublishOffer)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/offers/{id:[0-9]+}/image",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.UploadOfferImage)))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/images/{id:[0-9]+}",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.DeleteOfferImage)))).
		Methods(http.MethodDelete)

	// ЖК
	r.HandleFunc("/api/v1/zhk/{id:[0-9]+}", zhkHandler.GetZhkInfo).Methods("GET")
	r.HandleFunc("/api/v1/zhks", zhkHandler.GetAllZhk).Methods("GET")

	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)

	// Запуск сервера
	if err := http.ListenAndServe(":8001", corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
