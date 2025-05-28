package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetOffersHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{
		App: config.AppConfig{
			CORS: config.CORSConfig{AllowOrigin: "*"},
		},
	}
	offerHandlers := NewOfferHandler(mockUC, cfg)

	//t.Run("GetOffers ok", func(t *testing.T) {
	//	user := 1
	//	userPtr := &user
	//	ctx := context.WithValue(context.Background(), utils.UserIDKey, userPtr)
	//
	//	expectedOffers := []domain.OfferInfo{
	//		{Offer: domain.Offer{ID: 1}},
	//		{Offer: domain.Offer{ID: 2}},
	//	}
	//
	//	mockUC.EXPECT().GetOffers(gomock.Any(), userPtr).Return(expectedOffers, nil)
	//
	//	request := httptest.NewRequest(http.MethodGet, "/offers", nil).WithContext(ctx)
	//	response := httptest.NewRecorder()
	//
	//	offerHandlers.GetOffersHandler(response, request)
	//
	//	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	//})

	t.Run("GetOffersWithFilter ok", func(t *testing.T) {
		expectedOffers := []domain.OfferInfo{
			{Offer: domain.Offer{ID: 1}},
		}

		minPrice := 1000000
		mockUC.EXPECT().GetOffersByFilter(gomock.Any(), gomock.Eq(domain.OfferFilter{
			MinPrice: &minPrice,
		}), gomock.Any()).Return(expectedOffers, nil)

		request := httptest.NewRequest(http.MethodGet, "/offers?min_price=1000000", nil)
		response := httptest.NewRecorder()

		offerHandlers.GetOffersHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("GerOffers error", func(t *testing.T) {
		mockUC.EXPECT().GetOffers(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("filter error"))

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
			GetOffersByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
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

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{
		App: config.AppConfig{
			CORS: config.CORSConfig{AllowOrigin: "*"},
		},
	}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("successful get offer by id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.SoftUserIDKey, ptr(1))
		request := httptest.NewRequest(http.MethodGet, "/offers/1", nil).WithContext(ctx)
		request.Header.Set("X-Real-IP", "127.0.0.1")
		vars := map[string]string{"id": "1"}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		expectedOffer := domain.OfferInfo{Offer: domain.Offer{ID: 1, Price: 100}}

		mockUC.EXPECT().GetOfferByID(gomock.Any(), 1, "127.0.0.1", ptr(1)).Return(expectedOffer, nil)

		handler.GetOfferByID(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var result domain.OfferInfo
		err := json.NewDecoder(response.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, expectedOffer.Offer.ID, result.Offer.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/offers/abc", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "abc"})
		request = request.WithContext(context.WithValue(context.Background(), utils.SoftUserIDKey, ptr(1)))
		response := httptest.NewRecorder()

		handler.GetOfferByID(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		var resp map[string]string
		err := json.NewDecoder(response.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, "Некорректный ID", resp["error"])
	})

	t.Run("offer not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.SoftUserIDKey, ptr(1))

		request := httptest.NewRequest(http.MethodGet, "/offers/2", nil).WithContext(ctx)
		request.Header.Set("X-Real-IP", "127.0.0.1")
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 2, "127.0.0.1", ptr(1)).
			Return(domain.OfferInfo{}, fmt.Errorf("не найден"))

		handler.GetOfferByID(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)

		var errResp map[string]string
		err := json.NewDecoder(response.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Объявление не найдено", errResp["error"])
	})
}

func TestCreateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	offerHandlers := NewOfferHandler(mockUS, cfg)

	t.Run("CreateOffer ok", func(t *testing.T) {
		req := domain.Offer{
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
		req := domain.Offer{
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

	mockUS := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	offerHandlers := NewOfferHandler(mockUS, cfg)

	t.Run("UpdateOffer ok", func(t *testing.T) {
		req := domain.Offer{
			ID:       1,
			SellerID: 1,
			Price:    100,
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
		req := domain.Offer{
			ID:       1,
			SellerID: 1,
			Price:    100,
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

	mockUS := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{
		App: config.AppConfig{
			CORS: config.CORSConfig{AllowOrigin: "*"},
		},
	}
	offerHandlers := NewOfferHandler(mockUS, cfg)

	t.Run("DeleteOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		request := httptest.NewRequest(http.MethodPost, "/offers/1", nil).WithContext(ctx)
		vars := map[string]string{
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().CheckAccessToOffer(gomock.Any(), 1, 1).Return(nil)
		mockUS.EXPECT().DeleteOffer(gomock.Any(), 1).Return(nil)

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
		mockUS.EXPECT().DeleteOffer(gomock.Any(), 1).Return(fmt.Errorf("DeleteOffer usecase failed"))

		offerHandlers.DeleteOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}

func TestDeleteOfferImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{
		App: config.AppConfig{
			CORS: config.CORSConfig{AllowOrigin: "*"},
		},
	}
	offerHandlers := NewOfferHandler(mockUS, cfg)

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

func TestPublishOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("PublishOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 10)
		req := httptest.NewRequest(http.MethodPost, "/offers/1/publish", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().PublishOffer(gomock.Any(), 1, 10).Return(nil)

		handler.PublishOffer(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, "Объявление опубликовано", resp["message"])
	})

	t.Run("UserID not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/offers/1/publish", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.PublishOffer(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		var resp map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, "UserID not found", resp["error"])
	})

	t.Run("Invalid ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 10)
		req := httptest.NewRequest(http.MethodPost, "/offers/abc/publish", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.PublishOffer(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, "Некорректный ID", resp["error"])
	})

	t.Run("PublishOffer usecase error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 10)
		req := httptest.NewRequest(http.MethodPost, "/offers/1/publish", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().PublishOffer(gomock.Any(), 1, 10).Return(fmt.Errorf("ошибка публикации"))

		handler.PublishOffer(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, "ошибка публикации", resp["error"])
	})
}

func TestUploadOfferImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("Upload ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		buf := new(bytes.Buffer)
		_ = png.Encode(buf, img)
		imageBytes := buf.Bytes()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", "test.png")
		require.NoError(t, err)
		_, err = part.Write(imageBytes)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
			Return(domain.OfferInfo{Offer: domain.Offer{SellerID: 1}}, nil)

		mockUC.EXPECT().
			SaveOfferImage(gomock.Any(), 1, gomock.Any()).
			Return(int64(42), nil)

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		err = json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, float64(42), resp["image_id"])
	})

	t.Run("UserID not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid offer ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		req := httptest.NewRequest(http.MethodPost, "/offers/abc/images", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_, err := writer.CreateFormFile("image", "test.png")
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
			Return(domain.OfferInfo{Offer: domain.Offer{SellerID: 2}}, nil)

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("File not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
			Return(domain.OfferInfo{Offer: domain.Offer{SellerID: 1}}, nil)

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid file format", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", "test.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte("not an image"))
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
			Return(domain.OfferInfo{Offer: domain.Offer{SellerID: 1}}, nil)

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("SaveOfferImage failed", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)

		// Создаём валидное изображение
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		buf := new(bytes.Buffer)
		_ = png.Encode(buf, img)
		imageBytes := buf.Bytes()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", "test.png")
		require.NoError(t, err)
		_, err = part.Write(imageBytes)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/offers/1/images", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		mockUC.EXPECT().
			GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
			Return(domain.OfferInfo{Offer: domain.Offer{SellerID: 1}}, nil)

		mockUC.EXPECT().
			SaveOfferImage(gomock.Any(), 1, gomock.Any()).
			Return(int64(0), fmt.Errorf("fail"))

		handler.UploadOfferImage(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestLikeOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 42)
		reqBody := domain.LikeRequest{OfferId: 1}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/offers/like", bytes.NewBuffer(body)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		expected := domain.LikesStat{IsLiked: true, Amount: 10}
		mockUC.EXPECT().LikeOffer(gomock.Any(), domain.LikeRequest{OfferId: 1, UserId: 42}).Return(expected, nil)

		handler.LikeOffer(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result domain.LikesStat
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("UserID not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/offers/like", nil)
		rec := httptest.NewRecorder()

		handler.LikeOffer(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("bad request body", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 42)
		req := httptest.NewRequest(http.MethodPost, "/offers/like", bytes.NewBufferString("bad")).WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.LikeOffer(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 42)
		reqBody := domain.LikeRequest{OfferId: 1}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/offers/like", bytes.NewBuffer(body)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockUC.EXPECT().LikeOffer(gomock.Any(), domain.LikeRequest{OfferId: 1, UserId: 42}).Return(domain.LikesStat{}, fmt.Errorf("fail"))

		handler.LikeOffer(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestPromoteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("PromoteOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		mockUC.EXPECT().CheckType(gomock.Any(), 1).Return(true, nil)
		paymentResponse := &domain.CreatePaymentResponse{OfferId: 2, PaymentUri: "someURI"}
		mockUC.EXPECT().PromoteOffer(gomock.Any(), 2, 1).Return(paymentResponse, nil)

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("UserID not found", func(t *testing.T) {
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body))
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Incorrect ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/-2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "-2"})
		response := httptest.NewRecorder()

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		invalidBody := "{invalid json}"

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBufferString(invalidBody)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("No access to offer", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(errors.New("no access"))

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("Invalid promotion type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 99} // несуществующий тип
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		mockUC.EXPECT().CheckType(gomock.Any(), 99).Return(false, nil)

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Error checking promotion type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		mockUC.EXPECT().CheckType(gomock.Any(), 1).Return(false, errors.New("db error"))

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("Error promoting offer", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.CreatePaymentRequest{Type: 1}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", bytes.NewBuffer(body)).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		mockUC.EXPECT().CheckType(gomock.Any(), 1).Return(true, nil)
		mockUC.EXPECT().PromoteOffer(gomock.Any(), 2, 1).Return(nil, errors.New("payment error"))

		handler.PromoteOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestPromoteCheckOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("PromoteCheckOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		validateResult := true
		mockUC.EXPECT().ValidateOffer(gomock.Any(), 2, 1).Return(&validateResult, nil)
		paymentData := domain.CheckPaymentResponse{OfferId: 2, IsActive: true, IsPaid: true, Days: 30}
		mockUC.EXPECT().CheckPayment(gomock.Any(), 1).Return(&paymentData, nil)

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("UserID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Invalid offer ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/-2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "-2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid purchase ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "-1"})
		response := httptest.NewRecorder()

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("No access to offer", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(errors.New("no access"))

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("ValidateOffer error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		mockUC.EXPECT().ValidateOffer(gomock.Any(), 2, 1).Return(nil, errors.New("validation error"))

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("Invalid purchase validation", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		validateResult := false
		mockUC.EXPECT().ValidateOffer(gomock.Any(), 2, 1).Return(&validateResult, nil)

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("CheckPayment error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodPost, "/offers/2/promote", nil).WithContext(ctx)
		request = mux.SetURLVars(request, map[string]string{"id": "2", "purchaseId": "1"})
		response := httptest.NewRecorder()

		mockUC.EXPECT().CheckAccessToOffer(gomock.Any(), 2, 1).Return(nil)
		validateResult := true
		mockUC.EXPECT().ValidateOffer(gomock.Any(), 2, 1).Return(&validateResult, nil)
		mockUC.EXPECT().CheckPayment(gomock.Any(), 1).Return(nil, errors.New("payment error"))

		handler.PromoteCheckOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetFavorites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("GetFavorites OK without offer_type_id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodGet, "/offers/favorites", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		expected := []domain.OfferInfo{
			{Offer: domain.Offer{ID: 1, Price: 100}},
			{Offer: domain.Offer{ID: 2, Price: 200}},
		}

		mockUC.EXPECT().GetFavorites(gomock.Any(), 1, nil).Return(expected, nil)

		handler.GetFavorites(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var result []domain.OfferInfo
		err := json.NewDecoder(response.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("GetFavorites OK with offer_type_id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodGet, "/offers/favorites?offer_type_id=2", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		offerTypeID := 2
		expected := []domain.OfferInfo{
			{Offer: domain.Offer{ID: 3, Price: 300}},
		}

		mockUC.EXPECT().GetFavorites(gomock.Any(), 1, &offerTypeID).Return(expected, nil)

		handler.GetFavorites(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var result []domain.OfferInfo
		err := json.NewDecoder(response.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})


	t.Run("UserID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/offers/favorites", nil)
		response := httptest.NewRecorder()

		handler.GetFavorites(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var errResp map[string]string
		err := json.NewDecoder(response.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "UserID not found", errResp["error"])
	})

	t.Run("Usecase returns error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		request := httptest.NewRequest(http.MethodGet, "/offers/favorites", nil).WithContext(ctx)
		response := httptest.NewRecorder()

		mockUC.EXPECT().GetFavorites(gomock.Any(), 1, nil).Return(nil, errors.New("db error"))

		handler.GetFavorites(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)

		var errResp map[string]string
		err := json.NewDecoder(response.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Ошибка при получении избранных", errResp["error"])
	})
}

func TestFavoriteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewOfferHandler(mockUC, cfg)

	t.Run("FavoriteOffer ok", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.FavoriteRequest{OfferId: 2}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/favorite", bytes.NewBuffer(body)).WithContext(ctx)
		response := httptest.NewRecorder()

		stat := domain.FavoriteStat{IsFavorited: true, Amount: 1}
		req := domain.FavoriteRequest{UserId: 1, OfferId: 2}
		mockUC.EXPECT().FavoriteOffer(gomock.Any(), req).Return(stat, nil)

		handler.FavoriteOffer(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("UserID not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/offers/favorite", nil)
		response := httptest.NewRecorder()

		handler.FavoriteOffer(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		invalidBody := "{invalid json}"

		request := httptest.NewRequest(http.MethodPost, "/offers/favorite", bytes.NewBufferString(invalidBody)).WithContext(ctx)
		response := httptest.NewRecorder()

		handler.FavoriteOffer(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("FavotireOffer UC failed", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, 1)
		reqBody := domain.FavoriteRequest{OfferId: 2}
		body, _ := json.Marshal(reqBody)

		request := httptest.NewRequest(http.MethodPost, "/offers/favorite", bytes.NewBuffer(body)).WithContext(ctx)
		response := httptest.NewRecorder()

		stat := domain.FavoriteStat{IsFavorited: true, Amount: 1}
		req := domain.FavoriteRequest{UserId: 1, OfferId: 2}
		mockUC.EXPECT().FavoriteOffer(gomock.Any(), req).Return(stat, errors.New("some error"))

		handler.FavoriteOffer(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func ptr(i int) *int {
	return &i
}