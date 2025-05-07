package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

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

//func TestGetOfferByID(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUC := mocks.NewMockOfferUsecase(ctrl)
//	cfg := &config.Config{
//		App: config.AppConfig{
//			CORS: config.CORSConfig{AllowOrigin: "*"},
//		},
//	}
//	offerHandlers := NewOfferHandler(mockUC, cfg)
//
//	//t.Run("GetOfferByID ok", func(t *testing.T) {
//	//	userID := 1
//	//	ctx := context.WithValue(context.Background(), utils.UserIDKey, &userID)
//	//
//	//	expectedOffer := domain.OfferInfo{
//	//		Offer: domain.Offer{
//	//			ID:    1,
//	//			Price: 1000000,
//	//		},
//	//	}
//	//
//	//	mockUC.EXPECT().
//	//		GetOfferByID(gomock.Any(), 1, gomock.Any(), &userID).
//	//		Return(expectedOffer, nil)
//	//
//	//	request := httptest.NewRequest("GET", "/offers/1", nil).WithContext(ctx)
//	//	vars := map[string]string{
//	//		"id": "1",
//	//	}
//	//	request = mux.SetURLVars(request, vars)
//	//	response := httptest.NewRecorder()
//	//
//	//	offerHandlers.GetOfferByID(response, request)
//	//
//	//	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
//	//
//	//	var result domain.OfferInfo
//	//	err := json.NewDecoder(response.Body).Decode(&result)
//	//	assert.NoError(t, err)
//	//	assert.Equal(t, expectedOffer, result)
//	//})
//
//	t.Run("offer not found", func(t *testing.T) {
//		user := 1
//		userPtr := &user
//		ctx := context.WithValue(context.Background(), utils.UserIDKey, userPtr)
//
//		mockUC.EXPECT().GetOfferByID(gomock.Any(), 999, gomock.Any(), userPtr).
//			Return(domain.OfferInfo{}, fmt.Errorf("offer not found"))
//
//		request := httptest.NewRequest("GET", "/offers/999", nil).WithContext(ctx)
//		vars := map[string]string{
//			"id": "999",
//		}
//		request = mux.SetURLVars(request, vars)
//		response := httptest.NewRecorder()
//
//		offerHandlers.GetOfferByID(response, request)
//
//		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
//
//		var errResp map[string]string
//		err := json.NewDecoder(response.Body).Decode(&errResp)
//		assert.NoError(t, err)
//		assert.Equal(t, "Объявление не найдено", errResp["error"])
//	})
//
//	t.Run("Invalid ID", func(t *testing.T) {
//		request := httptest.NewRequest("GET", "/offers/-1", nil)
//		vars := map[string]string{
//			"id": "-1",
//		}
//		request = mux.SetURLVars(request, vars)
//		response := httptest.NewRecorder()
//
//		offerHandlers.GetOfferByID(response, request)
//
//		assert.Equal(t, http.StatusBadRequest, response.Code)
//
//		var errResp map[string]string
//		err := json.NewDecoder(response.Body).Decode(&errResp)
//		assert.NoError(t, err)
//		assert.Equal(t, "Некорректный ID", errResp["error"])
//	})
//}

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
