package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	handlers "github.com/go-park-mail-ru/2025_1_404/microservices/ai/delivery/http"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAIHandler_EvaluateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAIUsecase(ctrl)
	cfg := &config.Config{}
	handler := handlers.NewAIHandler(mockUC, cfg)

	t.Run("success", func(t *testing.T) {
		offer := domain.Offer{
			OfferType:     "sale",
			PurchaseType:  "mortgage",
			PropertyType:  "flat",
			Renovation:    "euro",
			Floor:         3,
			TotalFloors:   5,
			Rooms:         2,
			Address:       "Test Street",
			Area:          60,
			CeilingHeight: 3,
		}
		body, _ := json.Marshal(offer)

		expected := domain.EvaluationResult{
			MarketPrice:       domain.MarketPrice{Total: 6000000, PerSquareMeter: 100000},
			PossibleCostRange: domain.PossibleCostRange{Min: 5800000, Max: 6200000},
		}

		mockUC.EXPECT().
			EvaluateOffer(gomock.Any(), offer).
			Return(&expected, nil)

		req := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), utils.UserIDKey, 1))
		rec := httptest.NewRecorder()

		handler.EvaluateOffer(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var result domain.EvaluationResult
		_ = json.NewDecoder(rec.Body).Decode(&result)
		assert.Equal(t, expected, result)
	})

	t.Run("missing user id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/evaluate", nil)
		rec := httptest.NewRecorder()

		handler.EvaluateOffer(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader([]byte("bad json")))
		req = req.WithContext(context.WithValue(req.Context(), utils.UserIDKey, 1))
		rec := httptest.NewRecorder()

		handler.EvaluateOffer(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		offer := domain.Offer{OfferType: "sale"}
		body, _ := json.Marshal(offer)

		mockUC.EXPECT().EvaluateOffer(gomock.Any(), offer).Return(nil, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), utils.UserIDKey, 1))
		rec := httptest.NewRecorder()

		handler.EvaluateOffer(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
