package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-park-mail-ru/2025_1_404/config"
	deliveryOffer "github.com/go-park-mail-ru/2025_1_404/microservices/ai/delivery/http"
	usecaseOffer "github.com/go-park-mail-ru/2025_1_404/microservices/ai/usecase"
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

	// Инициализация подключения к Redis
	redisRepo, err := redis.New(&cfg.Redis, l)
	if err != nil {
		log.Fatalf("не удалось подключиться к Redis: %v", err)
	}

	//aiRepo := repoOffer.NewAIRepository(dbpool, l)
	aiUC := usecaseOffer.NewAIUsecase(redisRepo, l, cfg)
	aiHandler := deliveryOffer.NewAIHandler(aiUC, cfg)

	// Маршруты
	r := mux.NewRouter()

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(utils.NotFoundHandler)

	// Объявления
	r.Handle("/api/v1/evaluateOffer",
		middleware.AuthHandler(l, &cfg.App.CORS, middleware.CSRFMiddleware(l, cfg, http.HandlerFunc(aiHandler.EvaluateOffer)))).
		Methods(http.MethodPost)

	// Метрики
	metrics, reg := metrics.NewMetrics("ai")
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	r.Use(middleware.MetricsMiddleware(metrics))
	// AccessLog middleware
	logMux := middleware.AccessLog(l, r)
	// CORS middleware
	corsMux := middleware.CORSHandler(logMux, &cfg.App.CORS)

	log.Println("AI микросервис запущен")

	http.Handle("/metrics", promhttp.Handler())
	// Запуск сервера
	if err := http.ListenAndServe(cfg.App.Http.Port, corsMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
