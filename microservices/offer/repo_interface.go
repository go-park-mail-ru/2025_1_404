package offer

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_offer_repo.go -package=mocks
type OfferRepository interface {
	CreateOffer(ctx context.Context, offer repository.Offer) (int64, error)
	GetOfferByID(ctx context.Context, id int64) (repository.Offer, error)
	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]repository.Offer, error)
	GetAllOffers(ctx context.Context) ([]repository.Offer, error)
	GetOffersByFilter(ctx context.Context, f domain.OfferFilter) ([]repository.Offer, error)
	UpdateOffer(ctx context.Context, offer repository.Offer) error
	DeleteOffer(ctx context.Context, id int64) error
	CreateImageAndBindToOffer(ctx context.Context, offerID int, uuid string) (int64, error)
	UpdateOfferStatus(ctx context.Context, offerID int, statusID int) error
	GetOfferData(ctx context.Context, offer domain.Offer) (domain.OfferData, error)
	GetOfferImageWithUUID(ctx context.Context, imageID int64) (int64, string, error)
	DeleteOfferImage(ctx context.Context, imageID int64) error
}
