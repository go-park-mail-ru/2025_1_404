package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	mockFS "github.com/go-park-mail-ru/2025_1_404/internal/filestorage/mocks"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository/offer"
	mockRepo "github.com/go-park-mail-ru/2025_1_404/internal/usecase/offer/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("successful get offers", func(t *testing.T) {
		// Тестовые данные
		repoOffers := []repository.Offer{
			{ID: 1},
			{ID: 2},
		}

		domainOffers := []domain.Offer{
			{ID: 1},
			{ID: 2},
		}

		mockRepo.EXPECT().GetAllOffers(ctx).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar1.jpg"},
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar2.jpg"},
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		result, err := offerUsecase.GetOffers(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"avatar1.jpg", result[0].OfferData.Seller.Avatar)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"image1.jpg", result[0].OfferData.Images[0].Image)
	})

	t.Run("get offers from repository failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().GetAllOffers(ctx).Return([]repository.Offer{}, expectedErr)

		result, err := offerUsecase.GetOffers(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("prepare offer info failed", func(t *testing.T) {
		repoOffers := []repository.Offer{
			{ID: 1},
		}
		domainOffers := []domain.Offer{
			{ID: 1},
		}
		expectedErr := fmt.Errorf("offer data get failed")

		mockRepo.EXPECT().GetAllOffers(ctx).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0]).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffers(ctx)

		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOffersByFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	minPrice := 100000
	filter := domain.OfferFilter{
		MinPrice: &minPrice,
	}

	t.Run("GetOffersByFilter ok", func(t *testing.T) {
		// Тестовые данные
		repoOffers := []repository.Offer{
			{ID: 1, Price: 100},
			{ID: 2, Price: 10000000},
		}

		domainOffers := []domain.Offer{
			{ID: 1, Price: 100},
			{ID: 2, Price: 10000000},
		}

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar1.jpg"},
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar2.jpg"},
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"avatar1.jpg", result[0].OfferData.Seller.Avatar)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"image1.jpg", result[0].OfferData.Images[0].Image)
	})

	t.Run("filter error from repository", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter).Return(nil, expectedErr)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("prepare offer info failed", func(t *testing.T) {
		repoOffers := []repository.Offer{
			{ID: 1},
		}
		domainOffers := []domain.Offer{
			{ID: 1},
		}
		expectedErr := fmt.Errorf("offer data get failed")

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0]).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter)

		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOfferByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	testID := 123

	t.Run("successful get offer by ID", func(t *testing.T) {
		// Тестовые данные
		repoOffer := repository.Offer{
			ID:    int64(testID),
			Price: 100000,
		}

		expectedOfferData := domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar.jpg"},
			Images: []domain.OfferImage{{Image: "image1.jpg"}},
		}

		mockRepo.EXPECT().GetOfferByID(ctx, int64(testID)).Return(repoOffer, nil)
		mockRepo.EXPECT().GetOfferData(ctx, gomock.Any()).Return(expectedOfferData, nil)

		result, err := offerUsecase.GetOfferByID(ctx, testID)

		assert.NoError(t, err)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"avatar.jpg", result.OfferData.Seller.Avatar)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"image1.jpg", result.OfferData.Images[0].Image)
	})

	t.Run("offer not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("offer not found")

		mockRepo.EXPECT().GetOfferByID(ctx, int64(testID)).Return(repository.Offer{}, expectedErr)

		result, err := offerUsecase.GetOfferByID(ctx, testID)

		assert.Error(t, err)
		assert.Equal(t, domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("prepare offer info failed", func(t *testing.T) {
		repoOffer := repository.Offer{
			ID: int64(testID),
		}
		expectedErr := fmt.Errorf("offer data get failed")

		mockRepo.EXPECT().GetOfferByID(ctx, int64(testID)).Return(repoOffer, nil)
		mockRepo.EXPECT().GetOfferData(ctx, gomock.Any()).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOfferByID(ctx, testID)

		assert.Error(t, err)
		assert.Equal(t, domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOffersBySellerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	sellerID := 123

	t.Run("successful get offers by seller", func(t *testing.T) {
		// Тестовые данные
		repoOffers := []repository.Offer{
			{ID: 1},
			{ID: 2},
		}
		domainOffers := []domain.Offer{
			{ID: 1},
			{ID: 2},
		}

		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar1.jpg"},
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1]).Return(domain.OfferData{
			Seller: domain.OfferSeller{Avatar: "avatar2.jpg"},
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, utils.BasePath+utils.ImagesPath+"avatar1.jpg", result[0].OfferData.Seller.Avatar)
	})

	t.Run("no offers found for seller", func(t *testing.T) {
		// Репозиторий возвращает пустой список
		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return([]repository.Offer{}, nil)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID)

		// Проверяем результаты
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return(nil, expectedErr)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("prepare offers info failed", func(t *testing.T) {
		repoOffers := []repository.Offer{
			{ID: 1},
		}
		expectedErr := fmt.Errorf("offer data get failed")

		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, gomock.Any()).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID)

		// Проверяем результаты
		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCreateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	testOffer := domain.Offer{
		Price:    100000,
		SellerID: 123,
	}

	t.Run("CreateOffer ok", func(t *testing.T) {
		expectedRepoOffer := repository.Offer{
			Price:    100000,
			SellerID: 123,
			StatusID: 2,
		}

		mockRepo.EXPECT().CreateOffer(ctx, expectedRepoOffer).Return(int64(1), nil)

		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().CreateOffer(ctx, gomock.Any()).Return(int64(0), expectedErr)
		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, expectedErr, err)
	})
}

func TestDEleteOfferImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepo.NewMockofferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockFS := mockFS.NewMockFileStorage(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockFS)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")
	imageID := 123
	userID := 456
	offerID := int64(789)
	testUUID := "test-uuid-123"

	t.Run("DeleteOfferImage ok", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(offerID, testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, offerID).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(nil)
		mockFS.EXPECT().Delete(testUUID).Return(nil)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.NoError(t, err)
	})

	t.Run("image not found", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(0), "", fmt.Errorf("not found"))

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, "изображение не найдено", err.Error())
	})

	t.Run("offer not found", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(offerID, testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, offerID).Return(repository.Offer{}, fmt.Errorf("not found"))

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, "объявление не найдено", err.Error())
	})

	t.Run("user is not owner", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(offerID, testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, offerID).Return(repository.Offer{SellerID: 999}, nil)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, "нет доступа к удалению этого изображения", err.Error())
	})

	t.Run("failed to delete image relation", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(offerID, testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, offerID).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(fmt.Errorf("db error"))

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, "ошибка при удалении связи с изображением", err.Error())
	})

	t.Run("failed to delete physical file but still success", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(offerID, testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, offerID).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(nil)
		mockFS.EXPECT().Delete(testUUID).Return(fmt.Errorf("file not found"))

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.NoError(t, err)
	})
}
