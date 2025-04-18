package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_offer.go -package=mocks

type offerUsecase interface {
	GetOffers(ctx context.Context) ([]domain.OfferInfo, error)
	GetOffersByFilter(ctx context.Context, filter domain.OfferFilter) ([]domain.OfferInfo, error)
	GetOfferByID(ctx context.Context, id int) (domain.OfferInfo, error)
	GetOffersBySellerID(ctx context.Context, sellerID int) ([]domain.OfferInfo, error)
	CreateOffer(ctx context.Context, offer domain.Offer) (int, error)
	UpdateOffer(ctx context.Context, offer domain.Offer) error
	DeleteOffer(ctx context.Context, id int) error
	SaveOfferImage(ctx context.Context, offerID int, upload filestorage.FileUpload) (int64, error)
	PublishOffer(ctx context.Context, offerID int, userID int) error
	DeleteOfferImage(ctx context.Context, imageID int, userID int) error
	PrepareOfferInfo(ctx context.Context, offer domain.Offer) (domain.OfferInfo, error)
	PrepareOffersInfo(ctx context.Context, offers []domain.Offer) ([]domain.OfferInfo, error)
	CheckAccessToOffer(ctx context.Context, offerID int, userID int) error
}
