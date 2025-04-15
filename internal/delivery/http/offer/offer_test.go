package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/delivery/http/offer/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetOffersHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockofferUsecase(ctrl)
	offerHandlers := NewOfferHandler(mockUC)

	t.Run("GetOffers ok", func(t *testing.T) {
		expectedOffers := []domain.OfferInfo{
			{Offer: domain.Offer{ID: 1}},
			{Offer: domain.Offer{ID: 2}},
		}
		mockUC.EXPECT().GetOffers(gomock.Any()).Return(expectedOffers, nil)

		request := httptest.NewRequest(http.MethodGet, "/offers", nil)
		response := httptest.NewRecorder()

		offerHandlers.GetOffersHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("GetOffersWithFilter ok", func(t *testing.T) {
		expectedOffers := []domain.OfferInfo{
			{Offer: domain.Offer{ID: 1}},
		}

		minPrice := 1000000
		mockUC.EXPECT().GetOffersByFilter(gomock.Any(), gomock.Eq(domain.OfferFilter{
			MinPrice: &minPrice,
		})).Return(expectedOffers, nil)

		request := httptest.NewRequest(http.MethodGet, "/offers?min_price=1000000", nil)
		response := httptest.NewRecorder()

		offerHandlers.GetOffersHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("GerOffers error", func(t *testing.T) {
		mockUC.EXPECT().GetOffers(gomock.Any()).Return(nil, fmt.Errorf("filter error"))

		request := httptest.NewRequest(http.MethodGet, "/offers", nil)
		response := httptest.NewRecorder()

		offerHandlers.GetOffersHandler(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)

		var errResp map[string]string
		err := json.NewDecoder(response.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Ошибка при получении объявлений", errResp["error"])
	})

	t.Run("GetOffersByFilter error", func(t *testing.T) {
		mockUC.EXPECT().
			GetOffersByFilter(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("filter error"))

		request := httptest.NewRequest(http.MethodGet, "/offers?min_price=1000000", nil)
		response := httptest.NewRecorder()

		offerHandlers.GetOffersHandler(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)

		var errResp map[string]string
		err := json.NewDecoder(response.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Ошибка при фильтрации объявлений", errResp["error"])
	})
}

func TestGetOfferByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockofferUsecase(ctrl)
	offerHandlers := NewOfferHandler(mockUC)

	t.Run("GetOfferByID ok", func(t *testing.T) {
		expectedOffer := domain.OfferInfo{
			Offer: domain.Offer{
				ID:    1,
				Price: 1000000,
			},
		}

		mockUC.EXPECT().GetOfferByID(gomock.Any(), 1).Return(expectedOffer, nil)

		request := httptest.NewRequest("GET", "/offers/1", nil)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		offerHandlers.GetOfferByID(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		var result domain.OfferInfo
		err := json.NewDecoder(response.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, expectedOffer, result)
	})

	t.Run("offer not found", func(t *testing.T) {
        mockUC.EXPECT().GetOfferByID(gomock.Any(), 999).Return(domain.OfferInfo{}, fmt.Errorf("offer not found"))

        request := httptest.NewRequest("GET", "/offers/999", nil)
		vars := map[string]string{
			"id": "999",
		}
		request = mux.SetURLVars(request, vars)
        response := httptest.NewRecorder()

        offerHandlers.GetOfferByID(response, request)

        assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
        
        var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "Объявление не найдено", errResp["error"])
    })

	t.Run("Invalid ID", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/offers/-1", nil)
		vars := map[string]string{
			"id": "-1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		offerHandlers.GetOfferByID(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "Некорректный ID", errResp["error"])
	})
}

func TestCreateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockofferUsecase(ctrl)

	offerHandlers := NewOfferHandler(mockUS)

	t.Run("CreateOffer ok", func(t *testing.T) {
		req := domain.Offer {
			ID: 1,
		}

		expect := req
		expect.SellerID = 1

		body, _ := json.Marshal(req)
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewBuffer(body)).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().CreateOffer(gomock.Any(), expect).Return(1, nil)

		offerHandlers.CreateOffer(response, request)

		assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	})

	t.Run("userID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/offers", nil)
		response := httptest.NewRecorder()

		offerHandlers.CreateOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("bad request fields ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewBufferString("bad json")).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		offerHandlers.CreateOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("CreateOffer usecase failed", func(t *testing.T) {
		req := domain.Offer {
			ID: 1,
		}

		expect := req
		expect.SellerID = 1

		body, _ := json.Marshal(req)
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewBuffer(body)).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		mockUS.EXPECT().CreateOffer(gomock.Any(), expect).Return(1, fmt.Errorf("CreateOffer failed"))

		offerHandlers.CreateOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})

}

func TestUpdateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockofferUsecase(ctrl)

	offerHandlers := NewOfferHandler(mockUS)

	t.Run("UpdateOffer ok", func(t *testing.T) {
		req := domain.Offer {
			ID: 1,
			SellerID: 1,
			Price: 100,
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		body, _ := json.Marshal(req)

		request := httptest.NewRequest(http.MethodPost, "/offers/1", bytes.NewBuffer(body)).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(nil)
		mockUS.EXPECT().UpdateOffer(gomock.Any(), req).Return(nil)

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("userID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil)
		response := httptest.NewRecorder()

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "UserID not found", errResp["error"])
	})

	t.Run("incorrect offer id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/offers/-1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "-1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "Некорректный ID", errResp["error"])
	})

	t.Run("offer not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(fmt.Errorf("объявление не найдено"))

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "объявление не найдено", errResp["error"])
	})

	t.Run("alien userID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(fmt.Errorf("нет доступа к этому объявлению"))

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "нет доступа к этому объявлению", errResp["error"])
		
	})

	t.Run("UpdateOffer usecase failed", func(t *testing.T) {
		req := domain.Offer {
			ID: 1,
			SellerID: 1,
			Price: 100,
		}

		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		body, _ := json.Marshal(req)

		request := httptest.NewRequest(http.MethodPost, "/offers/1", bytes.NewBuffer(body)).WithContext(ctx)
		request.Header.Set("Content-Type", "application/json")
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(nil)
		mockUS.EXPECT().UpdateOffer(gomock.Any(), req).Return(fmt.Errorf("UpdateOffer failed"))

		offerHandlers.UpdateOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)

		var errResp map[string]string
        err := json.NewDecoder(response.Body).Decode(&errResp)
        assert.NoError(t, err)
        assert.Equal(t, "Ошибка при обновлении", errResp["error"])
	})
	
}

func TestDeleteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockofferUsecase(ctrl)

	offerHandlers := NewOfferHandler(mockUS)

	t.Run("DeleteOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
	
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(nil)
		mockUS.EXPECT().DeleteOffer(gomock.Any(),1).Return(nil)

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("userID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil)
		response := httptest.NewRecorder()

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("offer id empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("Offer not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
	
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(fmt.Errorf("объявление не найдено"))

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	})

	t.Run("DeleteOffer usecase failed", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
	
		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(nil)
		mockUS.EXPECT().DeleteOffer(gomock.Any(),1).Return(fmt.Errorf("DeleteOffer usecase failed"))

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}

func TestDeleteOfferImage(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUS := mocks.NewMockofferUsecase(ctrl)
    offerHandlers := NewOfferHandler(mockUS)

    t.Run("DeleteOfferImage ok", func(t *testing.T) {
        ctx := context.WithValue(context.Background(), utils.UserIDKey, 123)
        request := httptest.NewRequest(http.MethodDelete, "/offers/images/456", nil).WithContext(ctx)
        request = mux.SetURLVars(request, map[string]string{"id": "456"})
        response := httptest.NewRecorder()

        mockUS.EXPECT().DeleteOfferImage(gomock.Any(), 456, 123).Return(nil)

        offerHandlers.DeleteOfferImage(response, request)

        assert.Equal(t, http.StatusOK, response.Result().StatusCode)
        
        var responseBody map[string]string
        err := json.NewDecoder(response.Body).Decode(&responseBody)
        assert.NoError(t, err)
        assert.Equal(t, "Изображение удалено", responseBody["message"])
    })

    t.Run("UserID not found in context", func(t *testing.T) {
        request := httptest.NewRequest(http.MethodDelete, "/offers/images/456", nil)
        response := httptest.NewRecorder()

        offerHandlers.DeleteOfferImage(response, request)

        assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
    })

    t.Run("Invalid image ID in URL", func(t *testing.T) {
        ctx := context.WithValue(context.Background(), utils.UserIDKey, 123)
        request := httptest.NewRequest(http.MethodDelete, "/offers/images/invalid", nil).WithContext(ctx)
        request = mux.SetURLVars(request, map[string]string{"id": "invalid"})
        response := httptest.NewRecorder()

        offerHandlers.DeleteOfferImage(response, request)

        assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
    })

    t.Run("Usecase returns error", func(t *testing.T) {
        ctx := context.WithValue(context.Background(), utils.UserIDKey, 123)
        request := httptest.NewRequest(http.MethodDelete, "/offers/images/456", nil).WithContext(ctx)
        request = mux.SetURLVars(request, map[string]string{"id": "456"})
        response := httptest.NewRecorder()

        mockUS.EXPECT().DeleteOfferImage(gomock.Any(), 456, 123).Return(fmt.Errorf("image not found"))

        offerHandlers.DeleteOfferImage(response, request)

        assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
    })
}