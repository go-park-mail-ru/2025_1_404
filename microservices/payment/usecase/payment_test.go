package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	mock "github.com/go-park-mail-ru/2025_1_404/microservices/payment/mocks"
	yookassa "github.com/go-park-mail-ru/2025_1_404/pkg/api/yookassa"
	mockYookassa "github.com/go-park-mail-ru/2025_1_404/pkg/api/yookassa/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{App: config.AppConfig{BaseFrontendDir: "http://localhost:8000"}}
	mockRepo := mock.NewMockPaymentRepository(ctrl)
	mockYookassa := mockYookassa.NewMockYookassaRepo(ctrl)
	logger := logger.NewStub()

	usecase := NewPaymentUsecase(mockRepo, mockYookassa, logger, cfg)

	t.Run("valid types return true", func(t *testing.T) {
		validTypes := []int{1, 2, 3}
		for _, pt := range validTypes {
			assert.True(t, usecase.CheckType(pt), "paymentType %d should be valid", pt)
		}
	})

	t.Run("invalid types return false", func(t *testing.T) {
		invalidTypes := []int{0, -1, 4, 100}
		for _, pt := range invalidTypes {
			assert.False(t, usecase.CheckType(pt), "paymentType %d should be invalid", pt)
		}
	})
}

func TestValidatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPaymentRepository(ctrl)
	mockYookassa := mockYookassa.NewMockYookassaRepo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{}

	usecase := NewPaymentUsecase(mockRepo, mockYookassa, logger, cfg)

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "test-request-id")

	t.Run("valid payment and matching offerId", func(t *testing.T) {
		mockRepo.EXPECT().GetPaymentById(ctx, 1).Return(&domain.OfferPayment{
			Id:      1,
			OfferId: 10,
		}, nil)

		ok := usecase.ValidatePayment(ctx, 1, 10)
		assert.True(t, ok)
	})

	t.Run("valid payment but non-matching offerId", func(t *testing.T) {
		mockRepo.EXPECT().GetPaymentById(ctx, 2).Return(&domain.OfferPayment{
			Id:      2,
			OfferId: 20,
		}, nil)

		ok := usecase.ValidatePayment(ctx, 2, 99)
		assert.False(t, ok)
	})

	t.Run("payment not found", func(t *testing.T) {
		mockRepo.EXPECT().GetPaymentById(ctx, 3).Return(nil, errors.New("not found"))

		ok := usecase.ValidatePayment(ctx, 3, 10)
		assert.False(t, ok)
	})
}

func TestCheckPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPaymentRepository(ctrl)
	mockYookassa := mockYookassa.NewMockYookassaRepo(ctrl)
	logger := logger.NewStub()
	usecase := NewPaymentUsecase(mockRepo, mockYookassa, logger, nil)
	ctx := context.Background()

	t.Run("success update and deactivate", func(t *testing.T) {
		payment := &domain.OfferPayment{
			Id:         1,
			OfferId:    100,
			YookassaId: "yoo123",
			Type:       2,
			IsActive:   true,
			IsPaid:     false,
		}
		paymentResponse := &yookassa.CreatePaymentResponse{
			Status: "succeeded",
			Paid:   true,
		}

		mockRepo.EXPECT().GetPaymentById(ctx, payment.Id).Return(payment, nil)
		mockYookassa.EXPECT().GetPayment("yoo123").Return(paymentResponse, nil)
		mockRepo.EXPECT().UpdatePayment(ctx, gomock.AssignableToTypeOf(&domain.OfferPayment{})).Return(nil)
		mockRepo.EXPECT().DeactivateAllPaymentsByOfferId(ctx, payment.OfferId).Return(nil)

		result, err := usecase.CheckPayment(ctx, payment.Id)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsPaid)
		assert.Equal(t, 7, result.Days) // type=2 => 7 дней из GetPaymentPeriods()
	})

	t.Run("payment not active skips yookassa check", func(t *testing.T) {
		payment := &domain.OfferPayment{
			Id:       2,
			IsActive: false,
		}
		mockRepo.EXPECT().GetPaymentById(ctx, payment.Id).Return(payment, nil)
		// yookassaRepo.GetPayment не вызывается

		result, err := usecase.CheckPayment(ctx, payment.Id)
		require.NoError(t, err)
		assert.Equal(t, payment.Id, result.Id)
	})

	t.Run("yookassa get payment error", func(t *testing.T) {
		payment := &domain.OfferPayment{
			Id:         3,
			IsActive:   true,
			YookassaId: "fail",
		}
		mockRepo.EXPECT().GetPaymentById(ctx, payment.Id).Return(payment, nil)
		mockYookassa.EXPECT().GetPayment("fail").Return(nil, errors.New("yoo error"))

		result, err := usecase.CheckPayment(ctx, payment.Id)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("update payment error", func(t *testing.T) {
		payment := &domain.OfferPayment{
			Id:         4,
			IsActive:   true,
			YookassaId: "yooUpdate",
			Type:       1,
		}
		paymentResponse := &yookassa.CreatePaymentResponse{
			Status: "succeeded",
			Paid:   true,
		}

		mockRepo.EXPECT().GetPaymentById(ctx, payment.Id).Return(payment, nil)
		mockYookassa.EXPECT().GetPayment("yooUpdate").Return(paymentResponse, nil)
		mockRepo.EXPECT().UpdatePayment(ctx, gomock.AssignableToTypeOf(&domain.OfferPayment{})).Return(errors.New("update fail"))

		result, err := usecase.CheckPayment(ctx, payment.Id)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("deactivate payments error", func(t *testing.T) {
		payment := &domain.OfferPayment{
			Id:         5,
			IsActive:   true,
			YookassaId: "yooDeactivate",
			Type:       1,
			OfferId:    999,
		}
		paymentResponse := &yookassa.CreatePaymentResponse{
			Status: "succeeded",
			Paid:   true,
		}

		mockRepo.EXPECT().GetPaymentById(ctx, payment.Id).Return(payment, nil)
		mockYookassa.EXPECT().GetPayment("yooDeactivate").Return(paymentResponse, nil)
		mockRepo.EXPECT().UpdatePayment(ctx, gomock.AssignableToTypeOf(&domain.OfferPayment{})).Return(nil)
		mockRepo.EXPECT().DeactivateAllPaymentsByOfferId(ctx, payment.OfferId).Return(errors.New("deactivate fail"))

		result, err := usecase.CheckPayment(ctx, payment.Id)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPaymentRepository(ctrl)
	mockYooKassa := mockYookassa.NewMockYookassaRepo(ctrl)
	logger := logger.NewStub()
	cfg := &config.Config{App: config.AppConfig{BaseFrontendDir: "http://localhost:8080"}}

	usecase := NewPaymentUsecase(mockRepo, mockYooKassa, logger, cfg)
	ctx := context.Background()

	request := &domain.CreatePaymentRequest{
		OfferId: 123,
		Type:    2,
	}

	t.Run("success create payment", func(t *testing.T) {
		offerPayment := &domain.OfferPayment{
			Id:      1,
			OfferId: int(request.OfferId),
			Type:    request.Type,
		}
		yooResponse := &yookassa.CreatePaymentResponse{
			Id: "yoo-payment-1",
			Confirmation: yookassa.Confirmation{ReturnUri: "https://payment.confirmation.uri"},
		}

		// Моки
		mockRepo.EXPECT().
			CreatePaymentForOfferId(ctx, int(request.OfferId), request.Type).
			Return(offerPayment, nil)

		mockYooKassa.EXPECT().
			CreatePayment(
				2990,
				"Объявление №123. Продвижение 7 дней",
				"http://localhost:8080/offer/123/check/1",
			).
			Return(yooResponse, nil)

		mockRepo.EXPECT().
			UpdatePayment(ctx, gomock.AssignableToTypeOf(&domain.OfferPayment{})).
			DoAndReturn(func(ctx context.Context, payment *domain.OfferPayment) error {
				assert.Equal(t, "yoo-payment-1", payment.YookassaId)
				return nil
			})

		resp, err := usecase.CreatePayment(ctx, request)
	
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, request.OfferId, resp.OfferId)
	})

	t.Run("unknown payment type", func(t *testing.T) {
		badRequest := &domain.CreatePaymentRequest{
			OfferId: 1,
			Type:    99, // нет такого тарифа
		}

		resp, err := usecase.CreatePayment(ctx, badRequest)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "неизвестный тариф")
	})

	t.Run("create payment for offer failed", func(t *testing.T) {
		mockRepo.EXPECT().
			CreatePaymentForOfferId(ctx, int(request.OfferId), request.Type).
			Return(nil, errors.New("repo error"))

		resp, err := usecase.CreatePayment(ctx, request)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "repo error")
	})

	t.Run("yookassa create payment failed", func(t *testing.T) {
		offerPayment := &domain.OfferPayment{
			Id:      1,
			OfferId: int(request.OfferId),
			Type:    request.Type,
		}

		mockRepo.EXPECT().
			CreatePaymentForOfferId(ctx, int(request.OfferId), request.Type).
			Return(offerPayment, nil)

		mockYooKassa.EXPECT().
			CreatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("yookassa error"))

		resp, err := usecase.CreatePayment(ctx, request)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "yookassa error")
	})

	t.Run("update payment failed", func(t *testing.T) {
		offerPayment := &domain.OfferPayment{
			Id:      1,
			OfferId: int(request.OfferId),
			Type:    request.Type,
		}
		yooResponse := &yookassa.CreatePaymentResponse{
			Id: "yoo-payment-1",
		}

		mockRepo.EXPECT().
			CreatePaymentForOfferId(ctx, int(request.OfferId), request.Type).
			Return(offerPayment, nil)

		mockYooKassa.EXPECT().
			CreatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(yooResponse, nil)

		mockRepo.EXPECT().
			UpdatePayment(ctx, gomock.AssignableToTypeOf(&domain.OfferPayment{})).
			Return(errors.New("update error"))

		resp, err := usecase.CreatePayment(ctx, request)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "update error")
	})
}
