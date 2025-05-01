package offer

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_offer.go -package=mocks

type OfferUsecase interface {
	GetOffers(ctx context.Context) ([]domain.OfferInfo, error)
	GetOffersByFilter(ctx context.Context, filter domain.OfferFilter) ([]domain.OfferInfo, error)
	GetOfferByID(ctx context.Context, id int) (domain.OfferInfo, error)
	GetOffersBySellerID(ctx context.Context, sellerID int) ([]domain.OfferInfo, error)
	CreateOffer(ctx context.Context, offer domain.Offer) (int, error)
	UpdateOffer(ctx context.Context, offer domain.Offer) error
	DeleteOffer(ctx context.Context, id int) error
	SaveOfferImage(ctx context.Context, offerID int, upload s3.Upload) (int64, error)
	PublishOffer(ctx context.Context, offerID int, userID int) error
	DeleteOfferImage(ctx context.Context, imageID int, userID int) error
	PrepareOfferInfo(ctx context.Context, offer domain.Offer) (domain.OfferInfo, error)
	PrepareOffersInfo(ctx context.Context, offers []domain.Offer) ([]domain.OfferInfo, error)
	CheckAccessToOffer(ctx context.Context, offerID int, userID int) error
	GetOffersByZhkId (ctx context.Context, zhkId int) ([]domain.Offer, error)
	GetStations(ctx context.Context) ([]domain.Metro, error)
	LikeOffer(ctx context.Context,like domain.LikeRequest) (domain.LikesStat, error)
}
