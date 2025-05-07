package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func AuthHandler(log logger.Logger, cfg *config.CORSConfig, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.SendErrorResponse(w, "Учётные данные не предоставлены", http.StatusUnauthorized, cfg)
			log.Warn("jwt token was not found")
			return
		}

		claims, err := utils.ParseJWT(cookie.Value)
		if err != nil {
			utils.SendErrorResponse(w, "Неверный токен", http.StatusUnauthorized, cfg)
			log.Warn("invalid jwt token")
			return
		}

		userID := claims.UserID

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}
