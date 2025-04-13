package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetZhkInfoHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockZhkUsecase(ctrl)

	zhkHandlers := NewZhkHandler(mockUS)
	
	t.Run("GetZhkInfo ok", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhk/{id}", nil)
		vars := map[string]string {
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		zhk := domain.Zhk{ID: 1}
		mockUS.EXPECT().GetZhkByID(gomock.Any(), int64(1)).Return(zhk, nil)
		mockUS.EXPECT().GetZhkInfo(gomock.Any(), zhk).Return(domain.ZhkInfo{ID:1}, nil)

		zhkHandlers.GetZhkInfo(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("invalid zhk id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhk/{id}", nil)
		response := httptest.NewRecorder()

		zhkHandlers.GetZhkInfo(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("invalid zhk id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhk/{id}", nil)
		vars := map[string]string {
			"id": "string",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		zhkHandlers.GetZhkInfo(response, request)
		
		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("zhk not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhk/{id}", nil)
		vars := map[string]string {
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetZhkByID(gomock.Any(), int64(1)).Return(domain.Zhk{}, fmt.Errorf("zhk not found"))

		zhkHandlers.GetZhkInfo(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	})

	t.Run("usecase GetZhkInfo failed", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhk/{id}", nil)
		vars := map[string]string {
			"id": "1",
		}
		request = mux.SetURLVars(request, vars)
		response := httptest.NewRecorder()

		zhk := domain.Zhk{ID: 1}
		mockUS.EXPECT().GetZhkByID(gomock.Any(), int64(1)).Return(zhk, nil)
		mockUS.EXPECT().GetZhkInfo(gomock.Any(), zhk).Return(domain.ZhkInfo{}, fmt.Errorf("GetZhkInfo failed"))

		zhkHandlers.GetZhkInfo(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}

func TestGetAllZhkHandler (t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUS := mocks.NewMockZhkUsecase(ctrl)

	zhkHandlers := NewZhkHandler(mockUS)

	t.Run("GetAllZhk ok", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhks", nil)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetAllZhk(gomock.Any()).Return([]domain.ZhkInfo{domain.ZhkInfo{}}, nil)

		zhkHandlers.GetAllZhk(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("usecase GetAllZhk failed", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/zhks", nil)
		response := httptest.NewRecorder()

		mockUS.EXPECT().GetAllZhk(gomock.Any()).Return([]domain.ZhkInfo{}, fmt.Errorf("GetAllZhk failed"))

		zhkHandlers.GetAllZhk(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	})
}
