package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	deliveryZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/delivery/http"
	repoZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/repository"
	usecaseZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/usecase"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
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

	ctx := context.Background()

	// Логгер
	l, err := logger.NewZapLogger(cfg.App.Logger.Level)
	if err != nil {
		log.Fatalf("не удалось создать логгер: %v", err)
	}
	defer l.Close()

	// Инициализация подключения к БД
	dbpool, err := database.NewPool(&cfg.Postgres, ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	

	// Хранилище файлов
	basePath := "./internal/static/upload"

	// Подключаемся к offer grpc
	conn, err := grpc.NewClient(fmt.Sprint("offer", cfg.App.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("не удалось подключиться к offer grpc: %v", err)
		return
	}
	defer conn.Close()

	offerService := offerpb.NewOfferServiceClient(conn)

	zhkRepo := repoZhk.NewZhkRepository(dbpool, l)
	zhkUC := usecaseZhk.NewZhkUsecase(zhkRepo, l, cfg, offerService)
	zhkHandler := deliveryZhk.NewZhkHandler(zhkUC, cfg)

	// Маршруты
	r := mux.NewRouter()

	// Static
	r.PathPrefix("/images/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filestorage.ServeFile(w, r, basePath)
	}))

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	// ЖК
	r.HandleFunc("/api/v1/zhk/{id:[0-9]+}", zhkHandler.GetZhkInfo).Methods("GET")
	r.HandleFunc("/api/v1/zhks", zhkHandler.GetAllZhk).Methods("GET")

	// Метрики
	metrics, reg := metrics.NewMetrics("zhk")
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	r.Use(middleware.MetricsMiddleware(metrics))
	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)
	log.Println("Zhk микросервис запущен")

	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
