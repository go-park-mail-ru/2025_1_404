package offer

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_offer.go -package=mocks

type OfferUsecase interface {
	GetOffers(ctx context.Context, userID *int) ([]domain.OfferInfo, error)
	GetOffersByFilter(ctx context.Context, filter domain.OfferFilter, userID *int) ([]domain.OfferInfo, error)
	GetOfferByID(ctx context.Context, id int, ip string, userID *int) (domain.OfferInfo, error)
	GetOffersBySellerID(ctx context.Context, sellerID int, userID *int) ([]domain.OfferInfo, error)
	CreateOffer(ctx context.Context, offer domain.Offer) (int, error)
	UpdateOffer(ctx context.Context, offer domain.Offer) error
	DeleteOffer(ctx context.Context, id int) error
	SaveOfferImage(ctx context.Context, offerID int, upload s3.Upload) (int64, error)
	PublishOffer(ctx context.Context, offerID int, userID int) error
	DeleteOfferImage(ctx context.Context, imageID int, userID int) error
	PrepareOfferInfo(ctx context.Context, offer domain.Offer, userID *int) (domain.OfferInfo, error)
	PrepareOffersInfo(ctx context.Context, offers []domain.Offer, userID *int) ([]domain.OfferInfo, error)
	CheckAccessToOffer(ctx context.Context, offerID int, userID int) error
	GetOffersByZhkId(ctx context.Context, zhkId int) ([]domain.Offer, error)
	GetStations(ctx context.Context) ([]domain.Metro, error)
	LikeOffer(ctx context.Context, like domain.LikeRequest) (domain.LikesStat, error)
	GetFavorites(ctx context.Context, userID int, offerTypeID *int) ([]domain.OfferInfo, error)
	IsFavorite(ctx context.Context, userID, offerID int) (bool, error)
	FavoriteOffer(ctx context.Context, req domain.FavoriteRequest) (domain.FavoriteStat, error)
	PromoteOffer(ctx context.Context, offerID int, paymentType int) (*domain.CreatePaymentResponse, error)
	CheckType(ctx context.Context, paymentType int) (bool, error)
	ValidateOffer(ctx context.Context, offerID int, purchaseId int) (*bool, error)
	CheckPayment(ctx context.Context, paymentId int) (*domain.CheckPaymentResponse, error)
}
