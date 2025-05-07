package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	service "github.com/go-park-mail-ru/2025_1_404/microservices/auth/delivery/grpc"
	deliveryAuth "github.com/go-park-mail-ru/2025_1_404/microservices/auth/delivery/http"
	repoAuth "github.com/go-park-mail-ru/2025_1_404/microservices/auth/repository"
	usecaseAuth "github.com/go-park-mail-ru/2025_1_404/microservices/auth/usecase"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	// Конфиг
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

	authRepo := repoAuth.NewAuthRepository(dbpool, l)
	authUC := usecaseAuth.NewAuthUsecase(authRepo, l, s3repo, cfg)
	authHandler := deliveryAuth.NewAuthHandler(authUC, cfg)

	// Маршруты
	r := mux.NewRouter()

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
	r.Handle("/api/v1/auth/me", middleware.AuthHandler(l, &cfg.App.CORS, http.HandlerFunc(authHandler.Me))).
		Methods(http.MethodPost)
	r.Handle("/api/v1/users/update",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(authHandler.Update)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(authHandler.UploadImage)))).
		Methods(http.MethodPut)
	r.Handle("/api/v1/users/image",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(authHandler.DeleteImage)))).
		Methods(http.MethodDelete)
	r.Handle("/api/v1/users/csrf", middleware.AuthHandler(l, &cfg.App.CORS, http.HandlerFunc(authHandler.GetCSRFToken))).
		Methods(http.MethodGet)
	// Метрики
	metrics, reg := metrics.NewMetrics("auth")
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	r.Use(middleware.MetricsMiddleware(metrics))
	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)

	// Запуск grpc
	listen, err := net.Listen("tcp", cfg.App.Grpc.Port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer()

	authpb.RegisterAuthServiceServer(grpcServer, service.NewAuthService(authUC, l))

	go func() {
		log.Println("Auth grpc запущен")
		if err := grpcServer.Serve(listen); err != nil {
			log.Printf("failed to serve grpc: %v", err)
			return
		}
	}()

	log.Println("Auth микросервис запущен")

	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
