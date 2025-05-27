package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHandler_CreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockPaymentUsecase(ctrl)
	cfg := &config.Config{
		App: config.AppConfig{
			CORS: config.CORSConfig{AllowOrigin: "*"},
		},
	}
	handler := NewPaymentHandler(mockUC, cfg)

	t.Run("success", func(t *testing.T) {
		// Подготовка запроса и контекста с UserID
		reqBody := []byte(`{"offer_id": 123, "type": 2}`)
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		req := httptest.NewRequest(http.MethodPost, "/payment/create", bytes.NewReader(reqBody)).WithContext(ctx)
		w := httptest.NewRecorder()

		expectedResp := &domain.CreatePaymentResponse{
			OfferId:    123,
			PaymentUri: "https://payment.com/pay/123",
		}

		mockUC.EXPECT().CreatePayment(gomock.Any(), &domain.CreatePaymentRequest{
			OfferId: 123,
			Type:    2,
		}).Return(expectedResp, nil)

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var actualResp domain.CreatePaymentResponse
		err := json.NewDecoder(w.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Equal(t, *expectedResp, actualResp)
	})

	t.Run("missing user id", func(t *testing.T) {
		reqBody := []byte(`{"offer_id": 123, "type": 2}`)
		req := httptest.NewRequest(http.MethodPost, "/payment/create", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		req := httptest.NewRequest(http.MethodPost, "/payment/create", bytes.NewReader([]byte(`invalid-json`))).WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		reqBody := []byte(`{"offer_id": 123, "type": 2}`)
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		req := httptest.NewRequest(http.MethodPost, "/payment/create", bytes.NewReader(reqBody)).WithContext(ctx)
		w := httptest.NewRecorder()

		mockUC.EXPECT().CreatePayment(gomock.Any(), &domain.CreatePaymentRequest{
			OfferId: 123,
			Type:    2,
		}).Return(nil, errors.New("something went wrong"))

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
