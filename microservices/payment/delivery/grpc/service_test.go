package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/mocks"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	offerpb "github.com/go-park-mail-ru/2025_1_404/proto/offer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOfferById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	svc := service.NewOfferService(mockUC, logger.NewStub())

	offer := domain.Offer{
		ID:             1,
		SellerID:       2,
		OfferTypeID:    1,
		PropertyTypeID: 2,
		StatusID:       1,
		RenovationID:   1,
		Price:          5000000,
		Area:           50,
		Floor:          3,
		TotalFloors:    9,
		Rooms:          2,
		Flat:           12,
		CeilingHeight:  275,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockUC.EXPECT().
		GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
		Return(domain.OfferInfo{Offer: offer}, nil)

	resp, err := svc.GetOfferById(context.Background(), &offerpb.GetOfferRequest{Id: 1})
	assert.NoError(t, err)
	assert.Equal(t, int32(1), resp.Offer.Id)
}

func TestGetOfferById_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	svc := service.NewOfferService(mockUC, logger.NewStub())

	mockUC.EXPECT().
		GetOfferByID(gomock.Any(), 1, "", gomock.Nil()).
		Return(domain.OfferInfo{}, errors.New("not found"))

	_, err := svc.GetOfferById(context.Background(), &offerpb.GetOfferRequest{Id: 1})
	assert.Error(t, err)
}

func TestGetOffersByZhkId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	svc := service.NewOfferService(mockUC, logger.NewStub())

	offers := []domain.Offer{
		{ID: 1, SellerID: 2, OfferTypeID: 1, PropertyTypeID: 2, StatusID: 1, RenovationID: 1, Price: 4000000, Area: 44, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, SellerID: 3, OfferTypeID: 1, PropertyTypeID: 2, StatusID: 1, RenovationID: 1, Price: 4200000, Area: 46, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockUC.EXPECT().
		GetOffersByZhkId(gomock.Any(), 10).
		Return(offers, nil)

	resp, err := svc.GetOffersByZhkId(context.Background(), &offerpb.GetOffersByZhkRequest{ZhkId: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Offers, 2)
	assert.Equal(t, int32(1), resp.Offers[0].Id)
	assert.Equal(t, int32(2), resp.Offers[1].Id)
}

func TestGetOffersByZhkId_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOfferUsecase(ctrl)
	svc := service.NewOfferService(mockUC, logger.NewStub())

	mockUC.EXPECT().
		GetOffersByZhkId(gomock.Any(), 10).
		Return(nil, errors.New("fail"))

	_, err := svc.GetOffersByZhkId(context.Background(), &offerpb.GetOffersByZhkRequest{ZhkId: 10})
	assert.Error(t, err)
}
