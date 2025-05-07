package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/internal/metrics"
	"github.com/gorilla/mux"
)

func MetricsMiddleware(metrics *metrics.Metrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			wrappedRW := wrapResponseWriter(w)

			defer func() {
				duration := time.Since(start)
				method := r.Method

				var path string
				if route := mux.CurrentRoute(r); route != nil {
					path, _ = route.GetPathTemplate()
				} else {
					path = r.URL.Path
				}

				metrics.RecordRequest(method, path, wrappedRW.status)
				metrics.RecordRequestDuration(method, path, duration)

				if wrappedRW.status >= http.StatusBadRequest {
					metrics.RecordError(method, path, wrappedRW.status)
				}
			}()

			next.ServeHTTP(wrappedRW, r)
		})
	}
}
