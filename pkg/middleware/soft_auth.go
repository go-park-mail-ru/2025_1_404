package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func SoftAuthHandler(log logger.Logger, cfg *config.Config, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID *int = nil

		cookie, err := r.Cookie("token")
		if err == nil {
			claims, err := utils.ParseJWT(cookie.Value)
			if err == nil {
				id := claims.UserID
				userID = &id
			} else {
				log.Warn("invalid jwt token")
			}
		}

		ctx := context.WithValue(r.Context(), utils.SoftUserIDKey, userID)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}
