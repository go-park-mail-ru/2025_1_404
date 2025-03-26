package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Секретный ключ (пока просто строка)
var jwtSecret = []byte("supersecretkey")

// Claims Структура claims (что мы будем хранить в токене)
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT Генерация JWT токена
func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен живёт 24 часа

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT Парсинг и валидация JWT
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
