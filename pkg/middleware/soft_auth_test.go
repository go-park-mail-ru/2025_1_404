package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

func TestSoftAuthHandler_ValidToken(t *testing.T) {
	// Генерируем валидный JWT
	cookie, _ := utils.GenerateJWT(1)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Cookie", fmt.Sprintf(`token=%s`, cookie))
	rr := httptest.NewRecorder()

	// Заглушка логгера
	log := logger.NewStub()

	// Обработчик, который проверяет наличие userID в контексте
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(utils.SoftUserIDKey).(*int)
		if userID == nil || *userID != 1 {
			t.Errorf("Ожидался userID=1, получено: %v", userID)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := SoftAuthHandler(log, &config.Config{}, handler)
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}
}

func TestSoftAuthHandler_NoToken(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    rr := httptest.NewRecorder()
    log := logger.NewStub()
    
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value(utils.SoftUserIDKey).(*int)
        if userID != nil {
            t.Errorf("Ожидался userID=nil, получено: %v", userID)
        }
        w.WriteHeader(http.StatusOK)
    })
    
    middleware := SoftAuthHandler(log, &config.Config{}, handler)
    middleware.ServeHTTP(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("Ожидался статус 200, получен %d", rr.Code)
    }
}