package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	offerpb "github.com/go-park-mail-ru/2025_1_404/proto/offer"
	offerProtoMock "github.com/go-park-mail-ru/2025_1_404/proto/offer/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetZhkByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockZhkRepository(ctrl)
	mockLogger := logger.NewStub()
	cfg := &config.Config{}
	mockOfferService := offerProtoMock.NewMockOfferServiceClient(ctrl)

	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger, cfg, mockOfferService)
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

	mockRepo := mocks.NewMockZhkRepository(ctrl)
	mockLogger := logger.NewStub()
	minioPath := "http://localhost:9090"
	imagesPath := "/images/"
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: minioPath},
		App: config.AppConfig{BaseImagesPath: imagesPath},
	}
	mockOfferService := offerProtoMock.NewMockOfferServiceClient(ctrl)

	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger, cfg, mockOfferService)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	zhkID := int64(1)

	metroID := 4
	zhk := domain.Zhk{
		ID:          zhkID,
		Name:        "ЖК Тестовый",
		Address:     "ул. Тестовая, 1",
		Description: "Описание ЖК",
		Developer:   "Застройщик",
		Phone:       "+7 999 123 45 67",
		MetroStationId: &metroID,
	}

	header := domain.ZhkHeader{
		Images: []string{"img1.jpg", "img2.jpg"},
	}

	characteristics := domain.ZhkCharacteristics{
		Class: "Комфорт",
	}

	// []*offerpb.Offer {{Id: 1}, {Id: 2}}
	offers := &offerpb.GetOffersByZhkResponse{}

	t.Run("GetZhkById ok", func(t *testing.T) {
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(offers, nil)
		mockRepo.EXPECT().GetZhkMetro(ctx, int64(*zhk.MetroStationId)).Return(domain.ZhkMetro{Id: 4, Station: "Бауманская"}, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(characteristics, nil)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.NoError(t, err)
		assert.Equal(t, zhk.ID, result.ID)
		assert.Equal(t, minioPath+imagesPath+"img1.jpg", result.Header.Images[0])
	})

	t.Run("GetZhkById error", func(t *testing.T) {
		expectedErr := fmt.Errorf("ЖК с таким id не найден")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(domain.Zhk{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetOffers error", func(t *testing.T) {
		expectedErr := fmt.Errorf("Не удалось получить предложения у ЖК")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(nil, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkMetro error", func(t *testing.T) {
		expectedErr := fmt.Errorf("Не удалось получить метро ЖК")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(offers, nil)
		mockRepo.EXPECT().GetZhkMetro(ctx, int64(*zhk.MetroStationId)).Return(domain.ZhkMetro{}, expectedErr)
		
		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkHeader error", func(t *testing.T) {
		expectedErr := fmt.Errorf("Не удалось получить картинки ЖК")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(offers, nil)
		mockRepo.EXPECT().GetZhkMetro(ctx, int64(*zhk.MetroStationId)).Return(domain.ZhkMetro{Id: 4, Station: "Бауманская"}, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(domain.ZhkHeader{}, expectedErr)
		
		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

	t.Run("GetZhkByCharacteristics", func(t *testing.T) {
		expectedErr := fmt.Errorf("Не удалось получить класс ЖК")
		mockRepo.EXPECT().GetZhkByID(ctx, zhk.ID).Return(zhk, nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(offers, nil)
		mockRepo.EXPECT().GetZhkMetro(ctx, int64(*zhk.MetroStationId)).Return(domain.ZhkMetro{Id: 4, Station: "Бауманская"}, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk).Return(domain.ZhkCharacteristics{}, expectedErr)

		result, err := zhkUsecase.GetZhkInfo(ctx, zhk.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.ZhkInfo{}, result)
	})

}

func TestGetAllZhk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockZhkRepository(ctrl)
	mockLogger := logger.NewStub()
	minioPath := "http://localhost:9090"
	imagesPath := "/images/"
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: minioPath},
		App: config.AppConfig{BaseImagesPath: imagesPath},
	}
	mockOfferService := offerProtoMock.NewMockOfferServiceClient(ctrl)

	zhkUsecase := NewZhkUsecase(mockRepo, mockLogger, cfg, mockOfferService)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	zhkID := int64(1)

	metroID := 4
	zhk := []domain.Zhk{
		domain.Zhk{
			ID:          zhkID,
			Name:        "ЖК Тестовый",
			Address:     "ул. Тестовая, 1",
			Description: "Описание ЖК",
			Developer:   "Застройщик",
			Phone:       "+7 999 123 45 67",
			MetroStationId: &metroID,
		},
	}

	header := domain.ZhkHeader{
		Images: []string{"img1.jpg", "img2.jpg"},
	}

	characteristics := domain.ZhkCharacteristics{
		Class: "Комфорт",
	}

	// []*offerpb.Offer {{Id: 1}, {Id: 2}}
	offers := &offerpb.GetOffersByZhkResponse{}

	t.Run("GetAllZhk ok", func(t *testing.T) {
		mockRepo.EXPECT().GetAllZhk(ctx).Return(zhk, nil)

		mockRepo.EXPECT().GetZhkByID(ctx, zhk[0].ID).Return(zhk[0], nil)
		mockOfferService.EXPECT().GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(zhkID)}).Return(offers, nil)
		mockRepo.EXPECT().GetZhkMetro(ctx, int64(*zhk[0].MetroStationId)).Return(domain.ZhkMetro{Id: 4, Station: "Бауманская"}, nil)
		mockRepo.EXPECT().GetZhkHeader(ctx, zhk[0]).Return(header, nil)
		mockRepo.EXPECT().GetZhkCharacteristics(ctx, zhk[0]).Return(characteristics, nil)

		result, err := zhkUsecase.GetAllZhk(ctx)

		assert.NoError(t, err)
		assert.Equal(t, zhk[0].ID, result[0].ID)
		assert.Equal(t, minioPath+imagesPath+"img1.jpg", result[0].Header.Images[0])
	})

	t.Run("GetAllZhk failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("Не удалось получить список ЖК")
		mockRepo.EXPECT().GetAllZhk(ctx).Return(nil, expectedErr )

		result, err := zhkUsecase.GetAllZhk(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, []domain.ZhkInfo([]domain.ZhkInfo(nil)), result)
	})

}