package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	deliveryOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/delivery/http"
	repoOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
	usecaseOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/usecase"
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
	s3repo, err := s3.New(&cfg.Minio, l)
	if err != nil {
		log.Printf("не удалось подключиться к s3: %v", err)
		return
	}

	offerRepo := repoOffer.NewOfferRepository(dbpool, l)
	offerUC := usecaseOffer.NewOfferUsecase(offerRepo, l, s3repo, cfg)
	offerHandler := deliveryOffer.NewOfferHandler(offerUC, cfg)

	// Маршруты
	r := mux.NewRouter()

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

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

	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)

	log.Println("Offers микросервис запущен")

	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
