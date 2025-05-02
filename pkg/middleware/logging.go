package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func AccessLog(log logger.Logger, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedRW := wrapResponseWriter(w)

		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), utils.RequestIDKey, requestID)

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.WithFields(logger.LoggerFields{
					"requestID": requestID,
					"err":       err,
					"trace":     debug.Stack(),
				}).Error("recovered from panic")
			}
		}()

		nextHandler.ServeHTTP(wrappedRW, r.WithContext(ctx))

		log.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"status":    wrappedRW.status,
			"path":      r.URL.Path,
			"method":    r.Method,
			"duration":  time.Since(start),
		}).Info("request completed")
	})
}
