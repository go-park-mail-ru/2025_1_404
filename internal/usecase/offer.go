package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type OfferUsecase struct {
	repo   repository.Repository
	logger logger.Logger
	fs     filestorage.FileStorage
}

func NewOfferUsecase(repo repository.Repository, logger logger.Logger, fs filestorage.FileStorage) *OfferUsecase {
	return &OfferUsecase{repo: repo, logger: logger, fs: fs}
}

func (u *OfferUsecase) GetOffers(ctx context.Context) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	
	offers, err := u.repo.GetAllOffers(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: get all offers failed")
		return nil, err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"count":len(offers),}).Info("Offer usecase: offers fetched")

	offersDTO := mapOffers(offers)
	
	offersInfo, err := u.prepareOffersInfo(ctx, offersDTO)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *OfferUsecase) GetOffersByFilter(ctx context.Context, filter domain.OfferFilter) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rawOffers, err := u.repo.GetOffersByFilter(ctx, filter)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: filter offers failed")
		return nil, err
	}

	u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"count":len(rawOffers),}).Info("Offer usecase: offers filtered successfully")

	offersDTO := mapOffers(rawOffers)
	
	offersInfo, err := u.prepareOffersInfo(ctx, offersDTO)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *OfferUsecase) GetOfferByID(ctx context.Context, id int) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offer, err := u.repo.GetOfferByID(ctx, int64(id))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Error("Offer usecase: get offer by id failed")
		return domain.OfferInfo{}, err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": id}).Info("Offer usecase: offer fetched")

	offerDTO := mapOffer(offer)
	
	offerInfo, err := u.prepareOfferInfo(ctx, offerDTO)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, err
	}

	return offerInfo, nil
}

func (u *OfferUsecase) GetOffersBySellerID(ctx context.Context, sellerID int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetOffersBySellerID(ctx, int64(sellerID))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "seller_id": sellerID, "err": err.Error()}).Error("Offer usecase: get offers by seller failed")
		return nil, err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "seller_id": sellerID, "count": len(offers)}).Info("Offer usecase: offers by seller fetched")
	
	offersDTO := mapOffers(offers)
	
	offersInfo, err := u.prepareOffersInfo(ctx, offersDTO)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *OfferUsecase) CreateOffer(ctx context.Context, offer domain.Offer) (int, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offer.StatusID = 2

	repoOffer := unmapOffer(offer)
	id, err := u.repo.CreateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: create offer failed")
		return 0, err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": id}).Info("Offer usecase: offer created")
	return int(id), nil
}

func (u *OfferUsecase) UpdateOffer(ctx context.Context, offer domain.Offer) error {
	requestID := ctx.Value(utils.RequestIDKey)

	repoOffer := unmapOffer(offer)
	err := u.repo.UpdateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: update offer failed")
		return err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID}).Info("Offer usecase: offer updated")
	return nil
}

func (u *OfferUsecase) DeleteOffer(ctx context.Context, id int) error {
	err := u.repo.DeleteOffer(ctx, int64(id))
	requestID := ctx.Value(utils.RequestIDKey)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": id, "err": err.Error()}).Error("Offer usecase: delete offer failed")
		return err
	}
	u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": id}).Info("Offer usecase: offer deleted")
	return nil
}

func (u *OfferUsecase) SaveOfferImage(ctx context.Context, offerID int, upload filestorage.FileUpload) (int64, error) {
	err := u.fs.Add(upload)
	if err != nil {
		return 0, err
	}

	return u.repo.CreateImageAndBindToOffer(ctx, offerID, upload.Name)
}

func (u *OfferUsecase) PublishOffer(ctx context.Context, offerID int, userID int) error {
	offer, err := u.repo.GetOfferByID(ctx, int64(offerID))
	if err != nil {
		return fmt.Errorf("объявление не найдено")
	}
	if int(offer.SellerID) != userID {
		return fmt.Errorf("нет доступа к публикации этого объявления")
	}
	if offer.StatusID != 2 { // 2 = Черновик
		return fmt.Errorf("объявление уже активно или завершено")
	}

	// Проверка обязательных полей
	if offer.Price <= 0 || offer.Area <= 0 || offer.Floor <= 0 ||
		offer.TotalFloors <= 0 || offer.Rooms <= 0 || offer.PropertyTypeID == 0 ||
		offer.RenovationID == 0 || offer.OfferTypeID == 0 || offer.StatusID == 0 {
		return fmt.Errorf("не все обязательные поля заполнены")
	}

	return u.repo.UpdateOfferStatus(ctx, offerID, 1)
}

func (u *OfferUsecase) DeleteOfferImage(ctx context.Context, imageID int, userID int) error {
	offerID, uuid, err := u.repo.GetOfferImageWithUUID(ctx, int64(imageID))
	if err != nil {
		return fmt.Errorf("изображение не найдено")
	}

	offer, err := u.repo.GetOfferByID(ctx, offerID)
	if err != nil {
		return fmt.Errorf("объявление не найдено")
	}
	if int(offer.SellerID) != userID {
		return fmt.Errorf("нет доступа к удалению этого изображения")
	}

	err = u.repo.DeleteOfferImage(ctx, int64(imageID))
	if err != nil {
		return fmt.Errorf("ошибка при удалении связи с изображением")
	}

	// удаляем физически файл
	if err := u.fs.Delete(uuid); err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"image_id": imageID,
			"uuid":     uuid,
			"err":      err.Error(),
		}).Warn("Ошибка при удалении файла")
	}

	return nil
}

func (u *OfferUsecase) prepareOfferInfo (ctx context.Context, offer domain.Offer) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offerData, err := u.repo.GetOfferData(ctx, offer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID,"err":err.Error(),"offer_id": offer.ID,}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, fmt.Errorf("offer data get failed")
	}
	
	if offerData.Seller.Avatar != "" {
		offerData.Seller.Avatar = utils.BasePath + utils.ImagesPath + offerData.Seller.Avatar
	}

	for i, img := range offerData.Images {
		offerData.Images[i].Image = utils.BasePath + utils.ImagesPath + img.Image
	}


	offerInfo := domain.OfferInfo {
		Offer: offer,
		OfferData: offerData,
	}
	
	
	return offerInfo, nil
}

func (u *OfferUsecase) prepareOffersInfo (ctx context.Context, offers []domain.Offer) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offersInfo := make([]domain.OfferInfo, 0, len(offers))
	for _, offer := range offers {
		offerInfo, err := u.prepareOfferInfo(ctx, offer)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "offerID":offer.ID}).Error("Offer usecase: prepareOffersInfo failed")
			return []domain.OfferInfo{}, err
		}
		offersInfo = append(offersInfo, offerInfo)
	}
	return offersInfo, nil
}

func mapOffer(o repository.Offer) domain.Offer {
	return domain.Offer{
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
	}
}

func mapOffers(raw []repository.Offer) []domain.Offer {
	offers := make([]domain.Offer, 0, len(raw))
	for _, o := range raw {
		offers = append(offers, mapOffer(o))
	}
	return offers
}

func unmapOffer(o domain.Offer) repository.Offer {
	return repository.Offer{
		ID:             int64(o.ID),
		SellerID:       int64(o.SellerID),
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
	}
}
