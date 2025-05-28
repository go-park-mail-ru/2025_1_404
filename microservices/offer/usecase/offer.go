package usecase

import (
	"context"
	"fmt"
	paymentpb "github.com/go-park-mail-ru/2025_1_404/proto/payment"
	"html"
	"sort"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/api/yandex"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/redis"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/s3"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_404/proto/auth"
)

type offerUsecase struct {
	repo           offer.OfferRepository
	logger         logger.Logger
	s3Repo         s3.S3Repo
	yandexRepo     yandex.YandexRepo
	cfg            *config.Config
	authService    authpb.AuthServiceClient
	paymentService paymentpb.PaymentServiceClient
	redisRepo      redis.RedisRepo
}

func NewOfferUsecase(repo offer.OfferRepository, logger logger.Logger, s3Repo s3.S3Repo, cfg *config.Config, authService authpb.AuthServiceClient, paymentService paymentpb.PaymentServiceClient, redisRepo redis.RedisRepo, yandexRepo yandex.YandexRepo) *offerUsecase {
	return &offerUsecase{repo: repo, logger: logger, s3Repo: s3Repo, cfg: cfg, authService: authService, paymentService: paymentService, redisRepo: redisRepo, yandexRepo: yandexRepo}
}

func (u *offerUsecase) GetOffers(ctx context.Context, userID *int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetAllOffers(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get all offers failed")
		return nil, err
	}

	offersDTO := mapOffers(offers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO, userID)

	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) GetOffersByFilter(ctx context.Context, filter domain.OfferFilter, userID *int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rawOffers, err := u.repo.GetOffersByFilter(ctx, filter, userID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: filter offers failed")
		return nil, err
	}

	offersDTO := mapOffers(rawOffers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO, userID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offers data failed")
		return []domain.OfferInfo{}, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) GetOfferByID(ctx context.Context, id int, ip string, userID *int) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offer, err := u.repo.GetOfferByID(ctx, int64(id))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Error("Offer usecase: get offer by id failed")
		return domain.OfferInfo{}, err
	}

	err = u.addView(ctx, id, ip)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id, "err": err.Error()}).Error("Offer usecase: add view to offer failed")
	}

	offerDTO := mapOffer(offer)

	offerInfo, err := u.PrepareOfferInfo(ctx, offerDTO, userID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, err
	}

	if userID == nil || int64(*userID) != offer.SellerID {
		offerInfo.OfferData.OfferStat.Views = nil
	}

	return offerInfo, nil
}

func (u *offerUsecase) addView(ctx context.Context, offerId int, ip string) error {
	requestID := ctx.Value(utils.RequestIDKey)

	if ip == "" {
		return nil
	}

	key := fmt.Sprintf("view:%d:%s", offerId, ip)

	_, err := u.redisRepo.Get(ctx, key)
	if err == nil {
		return nil
	}

	if !u.redisRepo.IsNotFound(err) {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID}).Warn("Get view from redis failed")
		return err
	}

	if err := u.redisRepo.Set(ctx, key, "1", 10*time.Minute); err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID}).Error("failed to set view in redis")
		return err
	}

	if err := u.repo.IncrementView(ctx, offerId); err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID}).Error("failed to increment view")
		return err
	}

	return nil
}

func (u *offerUsecase) GetOffersBySellerID(ctx context.Context, sellerID int, userID *int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetOffersBySellerID(ctx, int64(sellerID))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "seller_id": sellerID, "err": err.Error()}).Error("Offer usecase: get offers by seller failed")
		return nil, err
	}

	offersDTO := mapOffers(offers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO, userID)

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

	if offer.Address == nil {
		return 0, fmt.Errorf("не указан адрес")
	}
	*offer.Address = html.EscapeString(*offer.Address)
	coords, err := u.yandexRepo.GetCoordinatesOfAddress(*offer.Address)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get coordinates failed")
		return 0, fmt.Errorf("не удалось получить координаты по адресу")
	}
	offer.Longitude = strconv.FormatFloat(coords.Longitude, 'f', -1, 64)
	offer.Latitude = strconv.FormatFloat(coords.Latitude, 'f', -1, 64)

	offer.StatusID = domain.OfferStatusDraft

	repoOffer := unmapOffer(offer)
	id, err := u.repo.CreateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: create offer failed")
		return 0, err
	}

	if offer.Price > 0 {
		err = u.repo.AddOrUpdatePriceHistory(ctx, id, offer.Price)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: price history update failed")
			return 0, err
		}
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

	offer.Latitude = existing.Latitude
	offer.Longitude = existing.Longitude

	if offer.Description != nil {
		escaped := html.EscapeString(*offer.Description)
		offer.Description = &escaped
	}

	if offer.Address == nil {
		return fmt.Errorf("не указан адрес")
	}

	*offer.Address = html.EscapeString(*offer.Address)
	if existing.Address != nil && *existing.Address != *offer.Address {
		coords, err := u.yandexRepo.GetCoordinatesOfAddress(*offer.Address)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: get coordinates failed")
			return fmt.Errorf("не удалось получить координаты по адресу")
		}
		offer.Longitude = strconv.FormatFloat(coords.Longitude, 'f', -1, 64)
		offer.Latitude = strconv.FormatFloat(coords.Latitude, 'f', -1, 64)
	}

	// Оставляем прежний статус
	offer.StatusID = existing.StatusID

	// Удаляем историю если изменился статус
	if offer.OfferTypeID != existing.OfferTypeID {
		err = u.repo.DeletePriceHistory(ctx, int64(offer.ID))
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: price history delete failed")
			return err
		}

		err = u.repo.AddOrUpdatePriceHistory(ctx, int64(offer.ID), offer.Price)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: price history update failed")
			return err
		}
	}

	// Обновляем цену если изменилась
	if offer.Price != existing.Price {
		err = u.repo.AddOrUpdatePriceHistory(ctx, int64(offer.ID), offer.Price)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: price history update failed")
			return err
		}
	}

	repoOffer := unmapOffer(offer)

	err = u.repo.UpdateOffer(ctx, repoOffer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offer.ID, "err": err.Error()}).Error("Offer usecase: update offer failed")
		return err
	}

	return nil
}

func (u *offerUsecase) DeleteOffer(ctx context.Context, id int) error {
	_ = u.repo.DeletePriceHistory(ctx, int64(id))

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

	if offer.Address == nil || *offer.Address == "" {
		return fmt.Errorf("не указан адрес")
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

func (u *offerUsecase) GetOffersByZhkId(ctx context.Context, zhkId int) ([]domain.Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offers, err := u.repo.GetOffersByZhkId(ctx, zhkId)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestId": requestID, "err": err.Error(), "zhkId": zhkId}).Error("Offer usecase: getOffersByZhkId failed")
		return []domain.Offer{}, err
	}

	return offers, nil
}

func (u *offerUsecase) GetStations(ctx context.Context) ([]domain.Metro, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	stations, err := u.repo.GetStations(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestId": requestID, "err": err.Error()}).Error("Offer usecase: getStations failed")
		return []domain.Metro{}, err
	}

	return stations, nil
}

func (u *offerUsecase) LikeOffer(ctx context.Context, like domain.LikeRequest) (domain.LikesStat, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var likeStat domain.LikesStat

	isLiked, err := u.repo.IsOfferLiked(ctx, like)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: isOfferLiked failed")
		return likeStat, err
	}

	if isLiked {
		err = u.repo.DeleteLike(ctx, like)
		likeStat.IsLiked = false
	} else {
		err = u.repo.CreateLike(ctx, like)
		likeStat.IsLiked = true
	}
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: toggle like failed")
		return likeStat, err
	}

	total, err := u.repo.GetLikeStat(ctx, like)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: getLikeAmount failed")
		return likeStat, err
	}

	likeStat.Amount = total

	return likeStat, nil
}

func (u *offerUsecase) GetFavorites(ctx context.Context, userID int, offerTypeID *int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rawOffers, err := u.repo.GetFavorites(ctx, int64(userID), offerTypeID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "userID": userID, "err": err.Error()}).Error("Offer usecase: get favorites failed")
		return nil, err
	}

	offersDTO := mapOffers(rawOffers)

	offersInfo, err := u.PrepareOffersInfo(ctx, offersDTO, &userID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Offer usecase: prepare favorites failed")
		return nil, err
	}

	return offersInfo, nil
}

func (u *offerUsecase) FavoriteOffer(ctx context.Context, req domain.FavoriteRequest) (domain.FavoriteStat, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var stat domain.FavoriteStat

	isFav, err := u.repo.IsFavorite(ctx, req.UserId, req.OfferId)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID, "err": err.Error(),
		}).Error("Offer usecase: isOfferFavorited failed")
		return stat, err
	}

	if isFav {
		err = u.repo.RemoveFavorite(ctx, req.UserId, req.OfferId)
		stat.IsFavorited = false
	} else {
		err = u.repo.AddFavorite(ctx, req.UserId, req.OfferId)
		stat.IsFavorited = true
	}
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID, "err": err.Error(),
		}).Error("Offer usecase: toggle favorite failed")
		return stat, err
	}

	total, err := u.repo.GetFavoriteStat(ctx, req)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID, "err": err.Error(),
		}).Error("Offer usecase: getFavoriteAmount failed")
		return stat, err
	}

	stat.Amount = total
	return stat, nil
}

func (u *offerUsecase) IsFavorite(ctx context.Context, userID, offerID int) (bool, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	isFav, err := u.repo.IsFavorite(ctx, userID, offerID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "user_id": userID, "offer_id": offerID, "err": err.Error()}).Error("Offer usecase: is favourite check failed")
	}
	return isFav, err
}

func (u *offerUsecase) CheckType(ctx context.Context, paymentType int) (bool, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	checkTypeResponse, err := u.paymentService.CheckType(ctx, &paymentpb.CheckTypeRequest{Type: int32(paymentType)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("Offer usecase: check type failed")
		return false, err
	}
	return checkTypeResponse.IsValid, nil
}

func (u *offerUsecase) PromoteOffer(ctx context.Context, offerID int, paymentType int) (*domain.CreatePaymentResponse, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	createPaymentResponse, err := u.paymentService.CreatePayment(ctx, &paymentpb.CreatePaymentRequest{Type: int32(paymentType), OfferId: int32(offerID)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("Offer usecase: create payment failed")
		return nil, err
	}
	return &domain.CreatePaymentResponse{
		OfferId:    createPaymentResponse.OfferId,
		PaymentUri: createPaymentResponse.RedirectUri,
	}, err
}

func (u *offerUsecase) ValidateOffer(ctx context.Context, offerID int, purchaseId int) (*bool, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	validateResponse, err := u.paymentService.ValidatePayment(ctx, &paymentpb.ValidatePaymentRequest{PaymentId: int32(purchaseId), OfferId: int32(offerID)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("Offer usecase: validate payment failed")
		return nil, err
	}
	return &validateResponse.IsValid, err
}

func (u *offerUsecase) CheckPayment(ctx context.Context, paymentId int) (*domain.CheckPaymentResponse, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	checkPaymentResponse, err := u.paymentService.CheckPayment(ctx, &paymentpb.CheckPaymentRequest{PaymentId: int32(paymentId)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("Offer usecase: check payment failed")
		return nil, err
	}

	if checkPaymentResponse.IsActive && checkPaymentResponse.IsPaid {
		if err := u.repo.SetPromotesUntil(ctx, int(checkPaymentResponse.OfferId), time.Now().Add(time.Duration(checkPaymentResponse.Days)*24*time.Hour)); err != nil {
			u.logger.Error("failed to set promotes_until")
		}
	}

	return &domain.CheckPaymentResponse{
		OfferId:  int(checkPaymentResponse.OfferId),
		IsActive: checkPaymentResponse.IsActive,
		IsPaid:   checkPaymentResponse.IsPaid,
		Days:     int(checkPaymentResponse.Days),
	}, nil
}

func (u *offerUsecase) VerifyOffer(ctx context.Context, offerID int) error {
	return u.repo.VerifyOffer(ctx, offerID)
}

func (u *offerUsecase) RejectOffer(ctx context.Context, offerID int, comment string) error {
	return u.repo.RejectOffer(ctx, offerID, comment)
}

func (u *offerUsecase) PrepareOfferInfo(ctx context.Context, offer domain.Offer, userID *int) (domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offerData, err := u.repo.GetOfferData(ctx, offer, userID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "offer_id": offer.ID}).Error("Offer usecase: get offer data failed")
		return domain.OfferInfo{}, fmt.Errorf("offer data get failed")
	}

	offerData.Promotion = nil
	if userID != nil && *userID == offer.SellerID {
		offerData.Promotion = &domain.OfferPromotion{
			IsPromoted:    offer.PromotesUntil != nil && offer.PromotesUntil.After(time.Now()),
			PromotedUntil: offer.PromotesUntil,
		}
	}
	offerData.PromotionScore = float32(offerData.OfferStat.LikesStat.Amount) * u.cfg.App.Promotion.LikeScore
	if offer.PromotesUntil != nil && offer.PromotesUntil.After(time.Now()) {
		offerData.PromotionScore += u.cfg.App.Promotion.PromotionScore
	}

	seller, err := u.authService.GetUserById(ctx, &authpb.GetUserRequest{Id: int32(offer.SellerID)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Warn("Offer usecase: get seller failed")
	}

	offerData.Seller = domain.OfferSeller{
		FirstName: seller.User.FirstName,
		LastName:  seller.User.LastName,
		Avatar:    seller.User.Image,
		CreatedAt: seller.User.CreatedAt.AsTime(),
	}

	for i, img := range offerData.Images {
		offerData.Images[i].Image = u.cfg.Minio.Path + u.cfg.Minio.OffersBucket + img.Image
	}

	priceHistory, err := u.repo.GetPriceHistory(ctx, int64(offer.ID), 5)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{
			"requestID": requestID, "offerID": offer.ID, "err": err.Error(),
		}).Warn("не удалось получить историю цен")
	} else {
		offerData.Prices = priceHistory
	}

	offerInfo := domain.OfferInfo{
		Offer:     offer,
		OfferData: offerData,
	}

	return offerInfo, nil
}

func (u *offerUsecase) PrepareOffersInfo(ctx context.Context, offers []domain.Offer, userID *int) ([]domain.OfferInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offersInfo := make([]domain.OfferInfo, 0, len(offers))
	for _, offer := range offers {
		offerInfo, err := u.PrepareOfferInfo(ctx, offer, userID)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "offerID": offer.ID}).Error("Offer usecase: prepareOffersInfo failed")
			return []domain.OfferInfo{}, err
		}
		offersInfo = append(offersInfo, offerInfo)
	}
	sort.Slice(offersInfo, func(i, j int) bool {
		return offersInfo[i].OfferData.PromotionScore > offersInfo[j].OfferData.PromotionScore
	})
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
		Longitude:      o.Longitude,
		Latitude:       o.Latitude,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
		PromotesUntil:  o.PromotesUntil,
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
		Longitude:      o.Longitude,
		Latitude:       o.Latitude,
	}
}
