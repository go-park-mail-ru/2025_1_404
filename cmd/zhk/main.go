package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	deliveryZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/delivery/http"
	repoZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/repository"
	usecaseZhk "github.com/go-park-mail-ru/2025_1_404/microservices/zhk/usecase"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
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
	l := logger.NewStub()
	//defer l.Close()

	// Хранилище файлов
	basePath := "./internal/static/upload"

	zhkRepo := repoZhk.NewZhkRepository(dbpool, l)
	zhkUC := usecaseZhk.NewZhkUsecase(zhkRepo, l, cfg)
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
