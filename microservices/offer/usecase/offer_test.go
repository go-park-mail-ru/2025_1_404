package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/mocks"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/api/yandex"
	yaMock "github.com/go-park-mail-ru/2025_1_404/pkg/api/yandex/mocks"
	redisMock "github.com/go-park-mail-ru/2025_1_404/pkg/database/redis/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	s3Mock "github.com/go-park-mail-ru/2025_1_404/pkg/database/s3/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
	authService "github.com/go-park-mail-ru/2025_1_404/proto/auth/mocks"
	paymentpb "github.com/go-park-mail-ru/2025_1_404/proto/payment"
	paymentService "github.com/go-park-mail-ru/2025_1_404/proto/payment/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	Path   = "http://localhost:9000"
	Bucket = "/offers/"
)

func TestGetOffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	var userID = 1
	UserID := &userID

	User1 := &authpb.GetUserResponse{User: &authpb.User{Id: 3, FirstName: "Ivan", LastName: "Ivanov", Image: "image1.png", CreatedAt: timestamppb.New(time.Now())}}
	User2 := &authpb.GetUserResponse{User: &authpb.User{Id: 1, FirstName: "Maksim", LastName: "Maksimov", Image: "image2.png", CreatedAt: timestamppb.New(time.Now())}}
	History1 := []domain.OfferPriceHistory{{Price: 123, Date: time.Now()}, {Price: 245, Date: time.Now()}}
	History2 := []domain.OfferPriceHistory{{Price: 789, Date: time.Now()}, {Price: 456, Date: time.Now()}}
	t.Run("successful get offers", func(t *testing.T) {
		// Тестовые данные
		repoOffers := []repository.Offer{
			{ID: 1, SellerID: 3},
			{ID: 2, SellerID: 1},
		}

		domainOffers := []domain.Offer{
			{ID: 1, SellerID: 3},
			{ID: 2, SellerID: 1},
		}

		mockRepo.EXPECT().GetAllOffers(ctx).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[0].SellerID)}).Return(User1, nil)
		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[1].SellerID)}).Return(User2, nil)

		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[0].ID), 5).Return(History1, nil)
		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[1].ID), 5).Return(History2, nil)
		result, err := offerUsecase.GetOffers(ctx, UserID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, Path+Bucket+"image1.jpg", result[0].OfferData.Images[0].Image)
	})

	t.Run("get offers from repository failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")
		mockRepo.EXPECT().GetAllOffers(ctx).Return([]repository.Offer{}, expectedErr)

		result, err := offerUsecase.GetOffers(ctx, UserID)

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
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0], UserID).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffers(ctx, UserID)

		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOffersByFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	var userID = 1
	UserID := &userID

	User1 := &authpb.GetUserResponse{User: &authpb.User{Id: 3, FirstName: "Ivan", LastName: "Ivanov", Image: "image1.png", CreatedAt: timestamppb.New(time.Now())}}
	User2 := &authpb.GetUserResponse{User: &authpb.User{Id: 1, FirstName: "Maksim", LastName: "Maksimov", Image: "image2.png", CreatedAt: timestamppb.New(time.Now())}}
	History1 := []domain.OfferPriceHistory{{Price: 123, Date: time.Now()}, {Price: 245, Date: time.Now()}}
	History2 := []domain.OfferPriceHistory{{Price: 789, Date: time.Now()}, {Price: 456, Date: time.Now()}}

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

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter, UserID).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[0].SellerID)}).Return(User1, nil)
		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[1].SellerID)}).Return(User2, nil)

		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[0].ID), 5).Return(History1, nil)
		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[1].ID), 5).Return(History2, nil)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter, UserID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, Path+Bucket+"image1.jpg", result[0].OfferData.Images[0].Image)
	})

	t.Run("filter error from repository", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter, UserID).Return(nil, expectedErr)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter, UserID)

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

		mockRepo.EXPECT().GetOffersByFilter(ctx, filter, UserID).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0], UserID).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffersByFilter(ctx, filter, UserID)

		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOfferByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	var userID = 1
	UserID := &userID

	User1 := &authpb.GetUserResponse{User: &authpb.User{Id: 1, FirstName: "Ivan", LastName: "Ivanov", Image: "image1.png", CreatedAt: timestamppb.New(time.Now())}}
	History1 := []domain.OfferPriceHistory{{Price: 123, Date: time.Now()}, {Price: 245, Date: time.Now()}}

	testID := 123
	IP := "123.123.123.123"
	key := fmt.Sprintf("view:%d:%s", testID, IP)

	t.Run("successful get offer by ID", func(t *testing.T) {
		// Тестовые данные
		repoOffer := repository.Offer{
			ID:       int64(testID),
			Price:    100000,
			SellerID: 1,
		}
		domainOffer := domain.Offer{
			ID:       testID,
			Price:    100000,
			SellerID: 1,
		}

		expectedOfferData := domain.OfferData{
			Images: []domain.OfferImage{{Image: "image1.jpg"}},
		}

		mockRepo.EXPECT().GetOfferByID(ctx, int64(testID)).Return(repoOffer, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffer, UserID).Return(expectedOfferData, nil)
		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffer.SellerID)}).Return(User1, nil)
		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffer.ID), 5).Return(History1, nil)
		mockRedis.EXPECT().Get(ctx, key).Return("", fmt.Errorf("some error"))
		mockRedis.EXPECT().IsNotFound(fmt.Errorf("some error")).Return(true)
		mockRedis.EXPECT().Set(ctx, key, "1", 10*time.Minute).Return(nil)
		mockRepo.EXPECT().IncrementView(gomock.Any(), testID).Return(nil)

		result, err := offerUsecase.GetOfferByID(ctx, testID, IP, UserID)

		assert.NoError(t, err)
		assert.Equal(t, Path+Bucket+"image1.jpg", result.OfferData.Images[0].Image)
	})

	t.Run("offer not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("offer not found")

		mockRepo.EXPECT().GetOfferByID(ctx, int64(testID)).Return(repository.Offer{}, expectedErr)

		result, err := offerUsecase.GetOfferByID(ctx, testID, IP, UserID)

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
		mockRepo.EXPECT().GetOfferData(ctx, gomock.Any(), UserID).Return(domain.OfferData{}, expectedErr)
		mockRedis.EXPECT().Get(ctx, key).Return("", nil)

		result, err := offerUsecase.GetOfferByID(ctx, testID, IP, UserID)

		assert.Error(t, err)
		assert.Equal(t, domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOffersBySellerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	var userID = 1
	UserID := &userID

	User1 := &authpb.GetUserResponse{User: &authpb.User{Id: 3, FirstName: "Ivan", LastName: "Ivanov", Image: "image1.png", CreatedAt: timestamppb.New(time.Now())}}
	User2 := &authpb.GetUserResponse{User: &authpb.User{Id: 1, FirstName: "Maksim", LastName: "Maksimov", Image: "image2.png", CreatedAt: timestamppb.New(time.Now())}}
	History1 := []domain.OfferPriceHistory{{Price: 123, Date: time.Now()}, {Price: 245, Date: time.Now()}}
	History2 := []domain.OfferPriceHistory{{Price: 789, Date: time.Now()}, {Price: 456, Date: time.Now()}}

	sellerID := 123

	t.Run("successful get offers by seller", func(t *testing.T) {
		// Тестовые данные
		repoOffers := []repository.Offer{
			{ID: 1, SellerID: 3},
			{ID: 2, SellerID: 1},
		}
		domainOffers := []domain.Offer{
			{ID: 1, SellerID: 3},
			{ID: 2, SellerID: 1},
		}

		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return(repoOffers, nil)
		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[0], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 1, Image: "image1.jpg"}},
		}, nil)

		mockRepo.EXPECT().GetOfferData(ctx, domainOffers[1], UserID).Return(domain.OfferData{
			Images: []domain.OfferImage{{ID: 2, Image: "image2.jpg"}},
		}, nil)

		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[0].SellerID)}).Return(User1, nil)
		mockAuthService.EXPECT().GetUserById(ctx, &authpb.GetUserRequest{Id: int32(domainOffers[1].SellerID)}).Return(User2, nil)

		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[0].ID), 5).Return(History1, nil)
		mockRepo.EXPECT().GetPriceHistory(ctx, int64(domainOffers[1].ID), 5).Return(History2, nil)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID, UserID)
		assert.NoError(t, err)
		assert.Equal(t, Path+Bucket+"image1.jpg", result[0].OfferData.Images[0].Image)
		assert.Len(t, result, 2)

	})

	t.Run("no offers found for seller", func(t *testing.T) {
		// Репозиторий возвращает пустой список
		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return([]repository.Offer{}, nil)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID, UserID)

		// Проверяем результаты
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("database error")

		mockRepo.EXPECT().GetOffersBySellerID(ctx, int64(sellerID)).Return(nil, expectedErr)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID, UserID)

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
		mockRepo.EXPECT().GetOfferData(ctx, gomock.Any(), UserID).Return(domain.OfferData{}, expectedErr)

		result, err := offerUsecase.GetOffersBySellerID(ctx, sellerID, UserID)

		// Проверяем результаты
		assert.Error(t, err)
		assert.Equal(t, []domain.OfferInfo{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCreateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	description := "Описание"
	address := "Москва, Улица Пушкина, Дом Кукушкина, квартира 1"
	testOffer := domain.Offer{
		ID:          1,
		Price:       100000,
		SellerID:    123,
		Description: &description,
		Address:     &address,
	}

	expectedRepoOffer := repository.Offer{
		ID:          1,
		Price:       100000,
		SellerID:    123,
		StatusID:    2,
		Description: &description,
		Address:     &address,
		Latitude:    "123.2",
		Longitude:   "123.1",
	}

	t.Run("CreateOffer ok", func(t *testing.T) {
		mockYa.EXPECT().GetCoordinatesOfAddress(*testOffer.Address).Return(&yandex.Coordinates{Latitude: 123.2, Longitude: 123.1}, nil)
		mockRepo.EXPECT().CreateOffer(ctx, expectedRepoOffer).Return(int64(1), nil)
		mockRepo.EXPECT().AddOrUpdatePriceHistory(ctx, int64(testOffer.ID), testOffer.Price).Return(nil)

		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("empty address", func(t *testing.T) {
		expectedErr := fmt.Errorf("не указан адрес")
		testOffer.Address = nil

		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("get coordinates failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось получить координаты по адресу")
		badAddress := "incorrect address"
		testOffer.Address = &badAddress

		mockYa.EXPECT().GetCoordinatesOfAddress(*testOffer.Address).Return(&yandex.Coordinates{Latitude: 1.1, Longitude: 2.2}, expectedErr)
		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("createRepo offer failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось создать оффер")

		mockYa.EXPECT().GetCoordinatesOfAddress(*testOffer.Address).Return(&yandex.Coordinates{Latitude: 123.2, Longitude: 123.1}, nil)
		mockRepo.EXPECT().CreateOffer(ctx, gomock.AssignableToTypeOf(repository.Offer{})).Return(int64(0), expectedErr)

		id, err := offerUsecase.CreateOffer(ctx, testOffer)

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, expectedErr, err)
	})
}

func TestDeleteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	offerId := 1

	t.Run("DeleteOffer success", func(t *testing.T) {
		mockRepo.EXPECT().DeletePriceHistory(ctx, int64(offerId)).Return(nil)
		mockRepo.EXPECT().DeleteOffer(ctx, int64(offerId)).Return(nil)

		err := offerUsecase.DeleteOffer(ctx, offerId)

		assert.NoError(t, err)
	})

	t.Run("DeleteOffer failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось удалить объявление")

		mockRepo.EXPECT().DeletePriceHistory(ctx, int64(offerId)).Return(nil)
		mockRepo.EXPECT().DeleteOffer(ctx, int64(offerId)).Return(expectedErr)

		err := offerUsecase.DeleteOffer(ctx, offerId)

		assert.Error(t, err, expectedErr)
	})
}

func TestSaveOfferImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	offerId := 1
	file := s3.Upload{
		Filename: "someImage.png",
		Bucket:   Bucket,
	}

	t.Run("SaveOfferImage success", func(t *testing.T) {
		mockS3.EXPECT().Put(ctx, file).Return(file.Filename, nil)
		mockRepo.EXPECT().CreateImageAndBindToOffer(ctx, offerId, file.Filename).Return(int64(1), nil)

		id, err := offerUsecase.SaveOfferImage(ctx, offerId, file)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("SaveOfferImage failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось сохранить фото")
		mockS3.EXPECT().Put(ctx, file).Return("", expectedErr)

		id, err := offerUsecase.SaveOfferImage(ctx, offerId, file)

		assert.Error(t, err, expectedErr)
		assert.Equal(t, int64(0), id)
	})
}

func TestPublishOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	offerID := 1
	userID := 1
	address := "Москва, Улица Пушкина, Дом Кукушкина, квартира 1"
	repoOffer := repository.Offer{
		ID:             int64(offerID),
		Price:          100000,
		SellerID:       1,
		StatusID:       2,
		Area:           99,
		Floor:          123,
		TotalFloors:    3,
		Rooms:          2,
		PropertyTypeID: 1,
		RenovationID:   1,
		OfferTypeID:    1,
		Address:        &address,
	}

	t.Run("PublishOffer success", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, nil)
		mockRepo.EXPECT().UpdateOfferStatus(ctx, offerID, 1).Return(nil)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.NoError(t, err)
	})

	t.Run("GetOfferById repo failed", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось получить объявление")
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, expectedErr)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.Error(t, err, expectedErr)
	})

	t.Run("Incorrect UserID", func(t *testing.T) {
		expectedErr := fmt.Errorf("нет доступа к публикации этого объявления")
		repoOffer.SellerID = 333
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, nil)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.Error(t, err, expectedErr)
	})

	t.Run("Offer is finished", func(t *testing.T) {
		expectedErr := fmt.Errorf("объявление уже активно или завершено")
		repoOffer.StatusID = 1
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, nil)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.Error(t, err, expectedErr)
	})

	t.Run("Incorrect field", func(t *testing.T) {
		expectedErr := fmt.Errorf("не все обязательные поля заполнены")
		repoOffer.OfferTypeID = 0
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, nil)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.Error(t, err, expectedErr)
	})

	t.Run("Empty Address", func(t *testing.T) {
		expectedErr := fmt.Errorf("не указан адрес")
		repoOffer.Address = nil
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repoOffer, nil)

		err := offerUsecase.PublishOffer(ctx, offerID, userID)

		assert.Error(t, err, expectedErr)
	})
}

func TestDeleteOfferImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	offerID := 1
	imageID := 1
	userID := 456
	testUUID := "test-uuid-123"

	t.Run("DeleteOfferImage ok", func(t *testing.T) {
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(nil)
		mockS3.EXPECT().Remove(ctx, "offers", testUUID).Return(nil)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.NoError(t, err)
	})

	t.Run("image not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("изображение не найдено")
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, expectedErr)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("offer not found", func(t *testing.T) {
		expectedErr := fmt.Errorf("объявление не найдено")
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repository.Offer{SellerID: int64(userID)}, expectedErr)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("user is not owner", func(t *testing.T) {
		expectedErr := fmt.Errorf("нет доступа к удалению этого изображения")
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repository.Offer{SellerID: int64(333)}, nil)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to delete image relation", func(t *testing.T) {
		expectedErr := fmt.Errorf("ошибка при удалении связи с изображением")
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(expectedErr)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to delete physical file but still success", func(t *testing.T) {
		expectedErr := fmt.Errorf("не удалось удалить из s3")
		mockRepo.EXPECT().GetOfferImageWithUUID(ctx, int64(imageID)).Return(int64(offerID), testUUID, nil)
		mockRepo.EXPECT().GetOfferByID(ctx, int64(offerID)).Return(repository.Offer{SellerID: int64(userID)}, nil)
		mockRepo.EXPECT().DeleteOfferImage(ctx, int64(imageID)).Return(nil)
		mockS3.EXPECT().Remove(ctx, "offers", testUUID).Return(expectedErr)

		err := offerUsecase.DeleteOfferImage(ctx, imageID, userID)

		assert.NoError(t, err)
	})
}

func TestLikeOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	like := domain.LikeRequest{
		OfferId: 1,
		UserId:  1,
	}

	t.Run("Offer already liked — should delete like", func(t *testing.T) {
		mockRepo.EXPECT().IsOfferLiked(ctx, like).Return(true, nil)
		mockRepo.EXPECT().DeleteLike(ctx, like).Return(nil)
		mockRepo.EXPECT().GetLikeStat(ctx, like).Return(5, nil)

		stat, err := offerUsecase.LikeOffer(ctx, like)

		assert.NoError(t, err)
		assert.Equal(t, false, stat.IsLiked)
		assert.Equal(t, 5, stat.Amount)
	})

	t.Run("Offer not liked — should create like", func(t *testing.T) {
		mockRepo.EXPECT().IsOfferLiked(ctx, like).Return(false, nil)
		mockRepo.EXPECT().CreateLike(ctx, like).Return(nil)
		mockRepo.EXPECT().GetLikeStat(ctx, like).Return(6, nil)

		stat, err := offerUsecase.LikeOffer(ctx, like)

		assert.NoError(t, err)
		assert.Equal(t, true, stat.IsLiked)
		assert.Equal(t, 6, stat.Amount)
	})

	t.Run("IsOfferLiked returns error", func(t *testing.T) {
		mockRepo.EXPECT().IsOfferLiked(ctx, like).Return(false, errors.New("db error"))

		stat, err := offerUsecase.LikeOffer(ctx, like)

		assert.Error(t, err)
		assert.Equal(t, 0, stat.Amount)
	})

	t.Run("CreateLike returns error", func(t *testing.T) {
		mockRepo.EXPECT().IsOfferLiked(ctx, like).Return(false, nil)
		mockRepo.EXPECT().CreateLike(ctx, like).Return(errors.New("create error"))

		stat, err := offerUsecase.LikeOffer(ctx, like)

		assert.Error(t, err)
		assert.Equal(t, 0, stat.Amount)
	})

	t.Run("GetLikeStat returns error", func(t *testing.T) {
		mockRepo.EXPECT().IsOfferLiked(ctx, like).Return(false, nil)
		mockRepo.EXPECT().CreateLike(ctx, like).Return(nil)
		mockRepo.EXPECT().GetLikeStat(ctx, like).Return(0, errors.New("stat error"))

		stat, err := offerUsecase.LikeOffer(ctx, like)

		assert.Error(t, err)
		assert.Equal(t, true, stat.IsLiked)
		assert.Equal(t, 0, stat.Amount)
	})
}

func TestCheckPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("CheckPayment ok", func(t *testing.T) {
		now := time.Now()
		expectedUntil := now.Add(time.Duration(30) * 24 * time.Hour)

		resp := paymentpb.CheckPaymentResponse{OfferId: 1, IsActive: true, IsPaid: true, Days: 30}
		mockPaymentService.EXPECT().CheckPayment(ctx, &paymentpb.CheckPaymentRequest{PaymentId: int32(1)}).Return(&resp, nil)
		mockRepo.EXPECT().
			SetPromotesUntil(ctx, 1, gomock.AssignableToTypeOf(time.Time{})).
			DoAndReturn(func(_ context.Context, _ int, actualUntil time.Time) error {
				assert.WithinDuration(t, expectedUntil, actualUntil, time.Second)
				return nil
			})

		paymentResponse, err := offerUsecase.CheckPayment(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, true, paymentResponse.IsActive)
	})

	t.Run("CheckPayment failed", func(t *testing.T) {
		resp := paymentpb.CheckPaymentResponse{OfferId: 1, IsActive: true, IsPaid: true, Days: 30}
		mockPaymentService.EXPECT().CheckPayment(ctx, &paymentpb.CheckPaymentRequest{PaymentId: int32(1)}).Return(&resp, errors.New("some error"))

		_, err := offerUsecase.CheckPayment(ctx, 1)

		assert.Error(t, err, "some error")
	})
}

func TestValidateOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("ValidateOffer ok", func(t *testing.T) {
		resp := paymentpb.ValidatePaymentResponse{IsValid: true}
		mockPaymentService.EXPECT().
			ValidatePayment(ctx, &paymentpb.ValidatePaymentRequest{
				PaymentId: int32(100),
				OfferId:   int32(1),
			}).Return(&resp, nil)

		result, err := offerUsecase.ValidateOffer(ctx, 1, 100)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, *result)
	})

	t.Run("ValidateOffer failed", func(t *testing.T) {
		mockPaymentService.EXPECT().
			ValidatePayment(ctx, &paymentpb.ValidatePaymentRequest{
				PaymentId: int32(100),
				OfferId:   int32(1),
			}).Return(nil, errors.New("validation failed"))

		result, err := offerUsecase.ValidateOffer(ctx, 1, 100)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPromoteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("PromoteOffer ok", func(t *testing.T) {
		resp := &paymentpb.CreatePaymentResponse{
			OfferId:     42,
			RedirectUri: "https://payment.example.com/redirect",
		}
		mockPaymentService.EXPECT().
			CreatePayment(ctx, &paymentpb.CreatePaymentRequest{
				Type:    int32(1),
				OfferId: int32(42),
			}).
			Return(resp, nil)

		result, err := offerUsecase.PromoteOffer(ctx, 42, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(42), result.OfferId)
		assert.Equal(t, "https://payment.example.com/redirect", result.PaymentUri)
	})

	t.Run("PromoteOffer failed", func(t *testing.T) {
		mockPaymentService.EXPECT().
			CreatePayment(ctx, &paymentpb.CreatePaymentRequest{
				Type:    int32(2),
				OfferId: int32(13),
			}).
			Return(nil, errors.New("create payment error"))

		result, err := offerUsecase.PromoteOffer(ctx, 13, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCheckType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("CheckType valid", func(t *testing.T) {
		resp := &paymentpb.CheckTypeResponse{IsValid: true}
		mockPaymentService.EXPECT().
			CheckType(ctx, &paymentpb.CheckTypeRequest{Type: int32(1)}).
			Return(resp, nil)

		isValid, err := offerUsecase.CheckType(ctx, 1)

		assert.NoError(t, err)
		assert.True(t, isValid)
	})

	t.Run("CheckType not valid", func(t *testing.T) {
		resp := &paymentpb.CheckTypeResponse{IsValid: false}
		mockPaymentService.EXPECT().
			CheckType(ctx, &paymentpb.CheckTypeRequest{Type: int32(2)}).
			Return(resp, nil)

		isValid, err := offerUsecase.CheckType(ctx, 2)

		assert.NoError(t, err)
		assert.False(t, isValid)
	})

	t.Run("CheckType failed", func(t *testing.T) {
		mockPaymentService.EXPECT().
			CheckType(ctx, &paymentpb.CheckTypeRequest{Type: int32(3)}).
			Return(nil, errors.New("check type error"))

		isValid, err := offerUsecase.CheckType(ctx, 3)

		assert.Error(t, err)
		assert.False(t, isValid)
	})
}

func TestIsFavorite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("IsFavorite true", func(t *testing.T) {
		mockRepo.EXPECT().
			IsFavorite(ctx, 42, 7).
			Return(true, nil)

		isFav, err := offerUsecase.IsFavorite(ctx, 42, 7)

		assert.NoError(t, err)
		assert.True(t, isFav)
	})

	t.Run("IsFavorite false", func(t *testing.T) {
		mockRepo.EXPECT().
			IsFavorite(ctx, 42, 8).
			Return(false, nil)

		isFav, err := offerUsecase.IsFavorite(ctx, 42, 8)

		assert.NoError(t, err)
		assert.False(t, isFav)
	})

	t.Run("IsFavorite error", func(t *testing.T) {
		mockRepo.EXPECT().
			IsFavorite(ctx, 42, 9).
			Return(false, errors.New("db error"))

		isFav, err := offerUsecase.IsFavorite(ctx, 42, 9)

		assert.Error(t, err)
		assert.False(t, isFav)
	})
}

func TestFavoriteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOfferRepository(ctrl)
	mockLogger := logger.NewStub()
	mockS3 := s3Mock.NewMockS3Repo(ctrl)
	cfg := &config.Config{
		Minio: config.MinioConfig{Path: Path, OffersBucket: Bucket},
	}
	mockYa := yaMock.NewMockYandexRepo(ctrl)
	mockPaymentService := paymentService.NewMockPaymentServiceClient(ctrl)
	mockAuthService := authService.NewMockAuthServiceClient(ctrl)
	mockRedis := redisMock.NewMockRedisRepo(ctrl)

	offerUsecase := NewOfferUsecase(mockRepo, mockLogger, mockS3, cfg, mockAuthService, mockPaymentService, mockRedis, mockYa)
	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	req := domain.FavoriteRequest{UserId: 1, OfferId: 10}

	t.Run("Add to favorites", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(false, nil)
		mockRepo.EXPECT().AddFavorite(ctx, req.UserId, req.OfferId).Return(nil)
		mockRepo.EXPECT().GetFavoriteStat(ctx, req).Return(5, nil)

		stat, err := offerUsecase.FavoriteOffer(ctx, req)

		assert.NoError(t, err)
		assert.True(t, stat.IsFavorited)
		assert.Equal(t, 5, stat.Amount)
	})

	t.Run("Remove from favorites", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(true, nil)
		mockRepo.EXPECT().RemoveFavorite(ctx, req.UserId, req.OfferId).Return(nil)
		mockRepo.EXPECT().GetFavoriteStat(ctx, req).Return(3, nil)

		stat, err := offerUsecase.FavoriteOffer(ctx, req)

		assert.NoError(t, err)
		assert.False(t, stat.IsFavorited)
		assert.Equal(t, 3, stat.Amount)
	})

	t.Run("Error on IsFavorite", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(false, errors.New("db error"))

		_, err := offerUsecase.FavoriteOffer(ctx, req)
		assert.Error(t, err)
	})

	t.Run("Error on AddFavorite", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(false, nil)
		mockRepo.EXPECT().AddFavorite(ctx, req.UserId, req.OfferId).Return(errors.New("add error"))

		_, err := offerUsecase.FavoriteOffer(ctx, req)
		assert.Error(t, err)
	})

	t.Run("Error on RemoveFavorite", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(true, nil)
		mockRepo.EXPECT().RemoveFavorite(ctx, req.UserId, req.OfferId).Return(errors.New("remove error"))

		_, err := offerUsecase.FavoriteOffer(ctx, req)
		assert.Error(t, err)
	})

	t.Run("Error on GetFavoriteStat", func(t *testing.T) {
		mockRepo.EXPECT().IsFavorite(ctx, req.UserId, req.OfferId).Return(true, nil)
		mockRepo.EXPECT().RemoveFavorite(ctx, req.UserId, req.OfferId).Return(nil)
		mockRepo.EXPECT().GetFavoriteStat(ctx, req).Return(0, errors.New("stat error"))

		_, err := offerUsecase.FavoriteOffer(ctx, req)
		assert.Error(t, err)
	})
}
