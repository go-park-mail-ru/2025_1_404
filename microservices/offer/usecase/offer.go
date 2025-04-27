package usecase

import (
	"context"
	"fmt"
	"html"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type offerUsecase struct {
	repo   offer.OfferRepository
	logger logger.Logger
	s3Repo s3.S3Repo
	cfg    *config.Config
}

func NewOfferUsecase(repo offer.OfferRepository, logger logger.Logger, s3Repo s3.S3Repo, cfg *config.Config) *offerUsecase {
	return &offerUsecase{repo: repo, logger: logger, s3Repo: s3Repo, cfg: cfg}
}

func (u *offerUsecase) GetOffers(ctx context.Context) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetAllOffers(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get all offers failed")
		return nil, err
	}

	offersDTO := mapOffers(offers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO)

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) GetOffersByFilter(ctx context.Context, filter domain.OfferFilter) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rawOffers, err := u.repo.GetOffersByFilter(ctx, filter)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: filter offers failed")
		return nil, err
	}

	offersDTO := mapOffers(rawOffers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) GetOfferByID(ctx context.Context, id int) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offer, err := u.repo.GetOfferByID(ctx, int64(id))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Error("Offer usecase: get offer by id failed")
		return domain.OfferInfo{}, err
	}

	offerDTO := mapOffer(offer)

	offerInfo, err := u.PrepareOfferInfo(ctx, offerDTO)

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, err
	}

	return offerInfo, nil
}

func (u *offerUsecase) GetOffersBySellerID(ctx context.Context, sellerID int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetOffersBySellerID(ctx, int64(sellerID))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "seller_id": sellerID, "err": err.Error()}).Error("Offer usecase: get offers by seller failed")
		return nil, err
	}

	offersDTO := mapOffers(offers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO)

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) CreateOffer(ctx context.Context, offer domain.Offer) (int, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	if offer.Description != nil {
		escaped := html.EscapeString(*offer.Description)
		offer.Description = &escaped
	}

	if offer.Address != nil {
		escaped := html.EscapeString(*offer.Address)
		offer.Address = &escaped
	}

	offer.StatusID = domain.OfferStatusDraft

	repoOffer := unmapOffer(offer)
	id, err := u.repo.CreateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: create offer failed")
		return 0, err
	}
	return int(id), nil
}

func (u *offerUsecase) UpdateOffer(ctx context.Context, offer domain.Offer) error {
	requestID := ctx.Value(utils.RequestIDKey)

	// Получаем существующее объявление
	existing, err := u.repo.GetOfferByID(ctx, int64(offer.ID))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: get offer before update failed")
		return fmt.Errorf("объявление не найдено")
	}

	if offer.Description != nil {
		escaped := html.EscapeString(*offer.Description)
		offer.Description = &escaped
	}

	if offer.Address != nil {
		escaped := html.EscapeString(*offer.Address)
		offer.Address = &escaped
	}

	// Оставляем прежний статус
	offer.StatusID = existing.StatusID

	repoOffer := unmapOffer(offer)

	err = u.repo.UpdateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: update offer failed")
		return err
	}

	return nil
}

func (u *offerUsecase) DeleteOffer(ctx context.Context, id int) error {
	err := u.repo.DeleteOffer(ctx, int64(id))
	requestID := ctx.Value(utils.RequestIDKey)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": id, "err": err.Error()}).Error("Offer usecase: delete offer failed")
		return err
	}
	return nil
}

func (u *offerUsecase) SaveOfferImage(ctx context.Context, offerID int, upload s3.Upload) (int64, error) {
	fileName, err := u.s3Repo.Put(ctx, upload)
	if err != nil {
		return 0, err
	}

	return u.repo.CreateImageAndBindToOffer(ctx, offerID, fileName)
}

func (u *offerUsecase) PublishOffer(ctx context.Context, offerID int, userID int) error {
	offer, err := u.repo.GetOfferByID(ctx, int64(offerID))
	if err != nil {
		return fmt.Errorf("объявление не найдено")
	}
	if int(offer.SellerID) != userID {
		return fmt.Errorf("нет доступа к публикации этого объявления")
	}
	if offer.StatusID != domain.OfferStatusDraft { // 2 = Черновик
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

func (u *offerUsecase) DeleteOfferImage(ctx context.Context, imageID int, userID int) error {
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
	if err := u.s3Repo.Remove(ctx, "offers", uuid); err != nil {
		u.logger.WithFields(logger.LoggerFields{"image_id": imageID, "uuid": uuid, "err": err.Error()}).Warn("Ошибка при удалении файла")
	}

	return nil
}

func (u *offerUsecase) PrepareOfferInfo(ctx context.Context, offer domain.Offer) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offerData, err := u.repo.GetOfferData(ctx, offer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "offer_id": offer.ID}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, fmt.Errorf("offer data get failed")
	}

	if offerData.Seller.Avatar != "" {
		offerData.Seller.Avatar = u.cfg.Minio.Path + u.cfg.Minio.AvatarsBucket + offerData.Seller.Avatar
	}

	for i, img := range offerData.Images {
		offerData.Images[i].Image = u.cfg.Minio.Path + u.cfg.Minio.OffersBucket + img.Image
	}

	offerInfo := domain.OfferInfo{
		Offer:     offer,
		OfferData: offerData,
	}

	return offerInfo, nil
}

func (u *offerUsecase) PrepareOffersInfo(ctx context.Context, offers []domain.Offer) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offersInfo := make([]domain.OfferInfo, 0, len(offers))
	for _, offer := range offers {
		offerInfo, err := u.PrepareOfferInfo(ctx, offer)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "offerID": offer.ID}).Error("Offer usecase: prepareOffersInfo failed")
			return []domain.OfferInfo{}, err
		}
		offersInfo = append(offersInfo, offerInfo)
	}
	return offersInfo, nil
}

func (u *offerUsecase) CheckAccessToOffer(ctx context.Context, offerID int, userID int) error {
	offer, err := u.repo.GetOfferByID(ctx, int64(offerID))
	if err != nil {
		return fmt.Errorf("объявление не найдено")
	}
	if int(offer.SellerID) != userID {
		return fmt.Errorf("нет доступа к этому объявлению")
	}
	return nil
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
