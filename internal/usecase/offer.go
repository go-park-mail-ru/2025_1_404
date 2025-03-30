package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
)

type OfferUsecase struct {
	repo repository.Repository
}

func NewOfferUsecase(repo repository.Repository) *OfferUsecase {
	return &OfferUsecase{repo: repo}
}

func (u *OfferUsecase) GetOffers(ctx context.Context) ([]domain.Offer, error) {
	repoOffers, err := u.repo.GetAllOffers(ctx)
	if err != nil {
		return nil, err
	}

	offers := make([]domain.Offer, 0, len(repoOffers))
	for _, o := range repoOffers {
		offers = append(offers, domain.Offer{
			ID:             int(o.ID),
			SellerID:       int(o.SellerID),
			OfferTypeID:    o.OfferTypeID,
			MetroStationID: o.MetroStationID,
			RentTypeID:     o.RentTypeID,
			PurchaseTypeID: o.PurchaseTypeID,
			PropertyTypeID: o.PropertyTypeID,
			StatusID:       o.StatusID,
			RenovationID:   o.RenovationID,
			ComplexID:      o.ComplexID,
			Price:          o.Price,
			Description:    o.Description,
			Floor:          o.Floor,
			TotalFloors:    o.TotalFloors,
			Rooms:          o.Rooms,
			Address:        o.Address,
			Flat:           o.Flat,
			Area:           o.Area,
			CeilingHeight:  o.CeilingHeight,
			CreatedAt:      o.CreatedAt,
			UpdatedAt:      o.UpdatedAt,
		})
	}

	return offers, nil
}
