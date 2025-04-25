package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/csrf"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func CSRFMiddleware(log logger.Logger, cfg *config.Config, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(cfg.App.Auth.CSRF.HeaderName)
		if token == "" {
			log.Warn("CSRF token header missing")
			utils.SendErrorResponse(w, "Необходим CSRF токен", http.StatusForbidden, &cfg.App.CORS)
			return
		}

		userID, ok := r.Context().Value(utils.UserIDKey).(int)
		if !ok {
			log.Error("userID not found in context")
			utils.SendErrorResponse(w, "user is not authorized", http.StatusForbidden, &cfg.App.CORS)
			return
		}

		if !csrf.ValidateCSRF(token, strconv.Itoa(userID), cfg.App.Auth.CSRF.Salt) {
			log.Error("incorrect CSRF token")
			utils.SendErrorResponse(w, "invalid csrf token", http.StatusForbidden, &cfg.App.CORS)
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}
