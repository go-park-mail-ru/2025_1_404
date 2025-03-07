package middleware

import (
	"github.com/go-park-mail-ru/2025_1_404/utils"
	"net/http"
)

func CORSHandler(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.EnableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}
