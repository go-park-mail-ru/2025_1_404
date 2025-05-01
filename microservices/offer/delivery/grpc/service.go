package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	offerpb "github.com/go-park-mail-ru/2025_1_404/proto/offer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type offerService struct {
	UC     offer.OfferUsecase
	logger logger.Logger
	offerpb.UnimplementedOfferServiceServer
}

func NewOfferService(usecase offer.OfferUsecase, logger logger.Logger) *offerService {
	return &offerService{UC: usecase, logger: logger, UnimplementedOfferServiceServer: offerpb.UnimplementedOfferServiceServer{}}
}

func (s *offerService) GetOfferById(ctx context.Context, r *offerpb.GetOfferRequest) (*offerpb.GetOfferResponse, error) {
	id := int(r.GetId())
	offer, err := s.UC.GetOfferByID(ctx, id, "", nil)
	if err != nil {
		s.logger.Warn("failed to get offer")
		return nil, status.Errorf(codes.NotFound, "cannot find offer by id: %v", err)
	}

	return &offerpb.GetOfferResponse{
		Offer: s.offerDtoToProto(&offer.Offer),
	}, nil
}

func (s *offerService) GetOffersByZhkId(ctx context.Context, r *offerpb.GetOffersByZhkRequest) (*offerpb.GetOffersByZhkResponse, error) {
	zhkId := int(r.GetZhkId())
	offers, err := s.UC.GetOffersByZhkId(ctx, zhkId)
	if err != nil {
		s.logger.Warn("failed to get zhk's offers")
		return nil, status.Errorf(codes.Internal, "failed to find offers by zhk id: %v", err)
	}

	return &offerpb.GetOffersByZhkResponse{
		Offers: s.offersDtoToProto(offers),
	}, nil
}

func (s *offerService) offersDtoToProto(offers []domain.Offer) []*offerpb.Offer {
	pbOffers := make([]*offerpb.Offer, 0, len(offers))
	for _, offer := range offers {
		pbOffers = append(pbOffers, s.offerDtoToProto(&offer))
	}
	return pbOffers
}

func (s *offerService) offerDtoToProto(offer *domain.Offer) *offerpb.Offer {
	return &offerpb.Offer{
		Id:             int32(offer.ID),
		SellerId:       int32(offer.SellerID),
		OfferTypeId:    int32(offer.OfferTypeID),
		MetroStationId: intToInt32Ptr(offer.MetroStationID),
		RentTypeId:     intToInt32Ptr(offer.RentTypeID),
		PurchaseTypeId: intToInt32Ptr(offer.PurchaseTypeID),
		PropertyTypeId: int32(offer.PropertyTypeID),
		StatusId:       int32(offer.StatusID),
		RenovationId:   int32(offer.RenovationID),
		ComplexId:      intToInt32Ptr(offer.ComplexID),
		Price:          int32(offer.Price),
		Description:    offer.Description,
		Floor:          int32(offer.Floor),
		TotalFloors:    int32(offer.TotalFloors),
		Rooms:          int32(offer.Rooms),
		Address:        offer.Address,
		Flat:           int32(offer.Flat),
		Area:           int32(offer.Area),
		Longitude:      offer.Longitude,
		Latitude:       offer.Latitude,
		CeilingHeight:  int32(offer.CeilingHeight),
		CreatedAt:      timestamppb.New(offer.CreatedAt),
		UpdatedAt:      timestamppb.New(offer.UpdatedAt),
	}
}

func intToInt32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}
