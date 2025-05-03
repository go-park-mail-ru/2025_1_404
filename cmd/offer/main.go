package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	service "github.com/go-park-mail-ru/2025_1_404/microservices/offer/delivery/grpc"
	deliveryOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/delivery/http"
	repoOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
	usecaseOffer "github.com/go-park-mail-ru/2025_1_404/microservices/offer/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/api/yandex"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/redis"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	offerpb "github.com/go-park-mail-ru/2025_1_404/proto/offer"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("не удалось загрузить конфиг: %v", err)
	}

	// Логгер
	l, err := logger.NewZapLogger(cfg.App.Logger.Level)
	if err != nil {
		log.Fatalf("не удалось создать логгер: %v", err)
	}
	defer l.Close()

	ctx := context.Background()


	// Инициализация подключения к БД
	dbpool, err := database.NewPool(&cfg.Postgres, ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	// Хранилище файлов
	s3repo, err := s3.New(&cfg.Minio, l)
	if err != nil {
		log.Printf("не удалось подключиться к s3: %v", err)
		return
	}

	yandexRepo := yandex.New(&cfg.Yandex)

	// Инициализация подключения к Redis
	redisRepo, err := redis.New(&cfg.Redis, l)
	if err != nil {
		log.Fatalf("не удалось подключиться к Redis: %v", err)
	}

	// Подключаемся к auth grpc
	conn, err := grpc.NewClient(fmt.Sprint("auth", cfg.App.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("не удалось подключиться к auth grpc: %v", err)
		return
	}
	defer conn.Close()

	authService := authpb.NewAuthServiceClient(conn)

	offerRepo := repoOffer.NewOfferRepository(dbpool, l)
	offerUC := usecaseOffer.NewOfferUsecase(offerRepo, l, s3repo, cfg, authService, redisRepo, yandexRepo)
	offerHandler := deliveryOffer.NewOfferHandler(offerUC, cfg)

	// Маршруты
	r := mux.NewRouter()

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	// Объявления
	r.Handle("/api/v1/offers",
		middleware.SoftAuthHandler(l, cfg, http.HandlerFunc(offerHandler.GetOffersHandler))).
		Methods(http.MethodGet)
	r.Handle("/api/v1/offers/{id:[0-9]+}",
		middleware.SoftAuthHandler(l, cfg, http.HandlerFunc(offerHandler.GetOfferByID))).
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
	r.HandleFunc("/api/v1/offers/stations", offerHandler.GetStations).
		Methods(http.MethodGet)
	r.Handle("/api/v1/offers/like",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(offerHandler.LikeOffer)))).
		Methods(http.MethodPost)

	// Метрики
	metrics, reg := metrics.NewMetrics("offer")
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	metricxMux := middleware.MetricsMiddleware(metrics, r)
	// AccessLog middleware
	logMux := middleware.AccessLog(l, metricxMux)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)

	// Запуск grpc
	listen, err := net.Listen("tcp", cfg.App.Grpc.Port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer()

	offerpb.RegisterOfferServiceServer(grpcServer, service.NewOfferService(offerUC, l))

	go func() {
		log.Println("Offer grpc запущен")
		if err := grpcServer.Serve(listen); err != nil {
			log.Printf("failed to serve grpc: %v", err)
			return
		}
	}()

	log.Println("Offers микросервис запущен")

	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
