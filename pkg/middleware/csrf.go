package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/pkg/csrf"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func CSRFMiddleware(log logger.Logger, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-CSRF-TOKEN")
		if token == "" {
			log.Warn("CSRF token header missing")
            utils.SendErrorResponse(w, "Необходим CSRF токен", http.StatusForbidden)
            return
		}

		userID, ok := r.Context().Value(utils.UserIDKey).(string)
		if !ok {
			log.Error("userID not found in context")
			utils.SendErrorResponse(w, "user is not authorizes", http.StatusForbidden)
			return
		}

		if !csrf.ValidateCSRF(token, userID, utils.Salt) {
			log.Error("incorrect CSRF token")
			utils.SendErrorResponse(w, "invalid csrf token", http.StatusForbidden)
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}
