package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var BasePath string
var ImagesPath string = "/images/"

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, переменные будут браться из окружения")
	}

	BasePath = os.Getenv("BASE_DIR")
	if BasePath == "" {
		BasePath = "http://localhost:8001"
	}
}
