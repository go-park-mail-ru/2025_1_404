package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	mockRepo "github.com/go-park-mail-ru/2025_1_404/internal/usecase/zhk/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetZhkByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockzhkRepository(ctrl)
	mockLogger := logger.NewStub()

	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	zhkID := int64(1)

	expectedZhk := domain.Zhk{
		ID:      zhkID,
		Name:    "ЖК Тестовый",
		Address: "ул. Тестовая, 1",
	}

	t.Run("GetZhkByID ok", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhkID).Return(expectedZhk, nil)

		result, err := zhkUsecase.GetZhkByID(ctx, zhkID)

		assert.NoError(t, err)
		assert.Equal(t, expectedZhk, result)
	})

	t.Run("zhk not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("zhk not found")

		mockRepo.EXPECT().GetZhkByID(ctx, zhkID).Return(domain.Zhk{}, expectedErr)

		result, err := zhkUsecase.GetZhkByID(ctx, zhkID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.Zhk{}, result)
	})
}

func TestGetZhkInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockzhkRepository(ctrl)
	mockLogger := logger.NewStub()
	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	zhk := domain.Zhk{
		ID:          1,
		Name:        "ЖК Тестовый",
		Address:     "ул. Тестовая, 1",
		Description: "Описание ЖК",
		Developer:   "Застройщик",
		Phone:       "+7 999 123 45 67",
	}

	header := domain.ZhkHeader{
		Images: []string{"img1.jpg", "img2.jpg"},
	}

	characteristics := domain.ZhkCharacteristics{
		Class: "Комфорт",
	}

	apartments := domain.ZhkApartments{
		Apartments: []domain.ZhkApartment{{Rooms: 1, LowestPrice: 5000000}},
	}

	reviews := domain.ZhkReviews{
		Reviews: []domain.Review{
			{Text: "Отлично", Avatar: "avatar.jpg"},
			{Text: "Норм", Avatar: ""},
		},
	}

	t.Run("GetZhkInfo ok", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(characteristics, nil)
		mockRepo.EXPECT().GetZhkApartments(ctx, zhk).Return(apartments, nil)
		mockRepo.EXPECT().GetZhkReviews(ctx, zhk).Return(reviews, nil)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.NoError(t, err)
		assert.Equal(t, zhk.ID, result.ID)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"img1.jpg", result.Header.Images[0])
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"avatar.jpg", result.Reviews.Reviews[0].Avatar)
		assert.Equal(t, "", result.Reviews.Reviews[1].Avatar)
	})

	t.Run("GetZhkHeader error", func(t *testing.T) {
		expectedErr := fmt.Errorf("header error")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(domain.ZhkHeader{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkCharacteristics error", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		expectedErr := fmt.Errorf("char error")
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(domain.ZhkCharacteristics{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkApartments error", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(characteristics, nil)
		expectedErr := fmt.Errorf("apartments error")
		mockRepo.EXPECT().GetZhkApartments(ctx, zhk).Return(domain.ZhkApartments{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkReviews error", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(characteristics, nil)
		mockRepo.EXPECT().GetZhkApartments(ctx, zhk).Return(apartments, nil)
		expectedErr := fmt.Errorf("reviews error")
		mockRepo.EXPECT().GetZhkReviews(ctx, zhk).Return(domain.ZhkReviews{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})
}
