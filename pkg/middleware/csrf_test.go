package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/csrf"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	l := logger.NewStub()
	cfg := &config.Config{
		App: config.AppConfig{ Auth: config.AuthConfig{ CSRF: config.CsrfStruct{Salt: "some_salt", HeaderName: "X-csrf-token"}}},
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	t.Run("CSRF ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPut, "/test", nil).WithContext(ctx)

		token := csrf.GenerateCSRF("1", cfg.App.Auth.CSRF.Salt)
		request.Header.Set("X-CSRF-TOKEN", token)

		response := httptest.NewRecorder()

		handler := CSRFMiddleware(l, cfg, testHandler)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusTeapot, response.Result().StatusCode)

	})

	t.Run("miss CSRF token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/test", nil)
		response := httptest.NewRecorder()

		handler := CSRFMiddleware(l, cfg, testHandler)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	})

	t.Run("context without userID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/test", nil)
		request.Header.Set("X-CSRF-TOKEN", "csrf_token")
		response := httptest.NewRecorder()

		handler := CSRFMiddleware(l, cfg, testHandler)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	})

	t.Run("incorrect CSRF token", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, "1")
		request := httptest.NewRequest(http.MethodPut, "/test", nil).WithContext(ctx)
		request.Header.Set("X-CSRF-TOKEN", "invalid_csrf_token")
		response := httptest.NewRecorder()

		handler := CSRFMiddleware(l, cfg, testHandler)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	})
}
