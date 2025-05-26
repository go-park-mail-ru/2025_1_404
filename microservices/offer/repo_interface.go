package offer

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_offer_repo.go -package=mocks
type OfferRepository interface {
	CreateOffer(ctx context.Context, offer repository.Offer) (int64, error)
	GetOfferByID(ctx context.Context, id int64) (repository.Offer, error)
	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]repository.Offer, error)
	GetAllOffers(ctx context.Context) ([]repository.Offer, error)
	GetOffersByFilter(ctx context.Context, f domain.OfferFilter, pUserId *int) ([]repository.Offer, error)
	UpdateOffer(ctx context.Context, offer repository.Offer) error
	DeleteOffer(ctx context.Context, id int64) error
	CreateImageAndBindToOffer(ctx context.Context, offerID int, uuid string) (int64, error)
	UpdateOfferStatus(ctx context.Context, offerID int, statusID int) error
	GetOfferData(ctx context.Context, offer domain.Offer, userID *int) (domain.OfferData, error)
	GetOfferImageWithUUID(ctx context.Context, imageID int64) (int64, string, error)
	DeleteOfferImage(ctx context.Context, imageID int64) error
	GetOffersByZhkId(ctx context.Context, zhkId int) ([]domain.Offer, error)
	GetStations(ctx context.Context) ([]domain.Metro, error)
	IsOfferLiked(ctx context.Context, like domain.LikeRequest) (bool, error)
	CreateLike(ctx context.Context, like domain.LikeRequest) error
	DeleteLike(ctx context.Context, like domain.LikeRequest) error
	GetLikeStat(ctx context.Context, like domain.LikeRequest) (int, error)
	IncrementView(ctx context.Context, id int) error
	AddOrUpdatePriceHistory(ctx context.Context, offerID int64, price int) error
	DeletePriceHistory(ctx context.Context, offerID int64) error
	GetPriceHistory(ctx context.Context, offerID int64, limit int) ([]domain.OfferPriceHistory, error)
	AddFavorite(ctx context.Context, userID, offerID int) error
	RemoveFavorite(ctx context.Context, userID, offerID int) error
	GetFavorites(ctx context.Context, userID int64, offerTypeID *int) ([]repository.Offer, error)
	IsFavorite(ctx context.Context, userID, offerID int) (bool, error)
	GetFavoriteStat(ctx context.Context, req domain.FavoriteRequest) (int, error)
	SetPromotesUntil(ctx context.Context, id int, until time.Time) error
}
