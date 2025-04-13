# --- Параметры по умолчанию ---
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_NAME ?= kvartirum
SSL_MODE ?= disable

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)

# Путь к проекту
ROOT_DIR := ../..
APP_DIR := $(ROOT_DIR)/app
MIGRATIONS_DIR := $(ROOT_DIR)/db/migrations

# Команды go
GO := go

# --- Сборка и запуск приложения ---
run:
	cd $(APP_DIR) && $(GO) run main.go

build:
	cd $(APP_DIR) && $(GO) build -o ../bin/app

# --- Миграции ---
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $$name