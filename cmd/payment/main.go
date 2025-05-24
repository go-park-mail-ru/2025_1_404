package main

import (
	"context"
	"fmt"
	service "github.com/go-park-mail-ru/2025_1_404/microservices/payment/delivery/grpc"
	repoPayment "github.com/go-park-mail-ru/2025_1_404/microservices/payment/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/api/yookassa"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	paymentpb "github.com/go-park-mail-ru/2025_1_404/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-park-mail-ru/2025_1_404/config"
	usecaseOffer "github.com/go-park-mail-ru/2025_1_404/microservices/payment/usecase"
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

	// Логгер
	l, err := logger.NewZapLogger(cfg.App.Logger.Level)
	if err != nil {
		log.Fatalf("не удалось создать логгер: %v", err)
	}
	defer l.Close()

	// Подключаемся к offer grpc
	conn, err := grpc.NewClient(fmt.Sprint("offer", cfg.App.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("не удалось подключиться к offer grpc: %v", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()
	// Инициализация подключения к БД
	dbpool, err := database.NewPool(&cfg.Postgres, ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	// Yookassa
	yookassaRepo := yookassa.New(&cfg.Yookassa)

	paymentRepo := repoPayment.NewPaymentRepository(dbpool, l)
	paymentUC := usecaseOffer.NewPaymentUsecase(paymentRepo, yookassaRepo, l, cfg)

	// Маршруты
	r := mux.NewRouter()

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	// Метрики
	metrics, reg := metrics.NewMetrics("payment")
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

	paymentpb.RegisterPaymentServiceServer(grpcServer, service.NewPaymentService(paymentUC, l))

	go func() {
		log.Println("Payment grpc запущен")
		if err := grpcServer.Serve(listen); err != nil {
			log.Printf("failed to serve grpc: %v", err)
			return
		}
	}()

	log.Println("Payment микросервис запущен")

	http.Handle("/metrics", promhttp.Handler())
	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
