package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func CORSHandler(nextHandler http.Handler, cfg *config.CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.EnableCORS(w, cfg)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}
