package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk"
	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	offerpb "github.com/go-park-mail-ru/2025_1_404/proto/offer"
)

type zhkUsecase struct {
	repo         zhk.ZhkRepository
	logger       logger.Logger
	cfg          *config.Config
	offerService offerpb.OfferServiceClient
}

func NewZhkUsecase(repo zhk.ZhkRepository, logger logger.Logger, cfg *config.Config, offerSerice offerpb.OfferServiceClient) *zhkUsecase {
	return &zhkUsecase{repo: repo, logger: logger, cfg: cfg, offerService: offerSerice}
}

func (u *zhkUsecase) GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	zhk, err := u.repo.GetZhkByID(ctx, id)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "id": id}).Warn("ЖК с таким id не существует")
		return domain.Zhk{}, err
	}
	return zhk, nil
}

func (u *zhkUsecase) GetZhkInfo(ctx context.Context, id int64) (domain.ZhkInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	// Получаем ЖК
	zhk, err := u.GetZhkByID(ctx, id)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "zhkId": id, "err": err.Error()}).Warn("failed to get zhk")
		return domain.ZhkInfo{}, errors.New("ЖК с таким id не найден")
	}

	// Получаем предложения ЖК
	zhkOffers, err := u.offerService.GetOffersByZhkId(ctx, &offerpb.GetOffersByZhkRequest{ZhkId: int32(id)})
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "zhkId": id, "err": err.Error()}).Warn("failed to get zhk offers")
		return domain.ZhkInfo{}, errors.New("Не удалось получить предложения у ЖК")
	}

	// Получаем информацию о метро у ЖК
	zhkMetro, err := u.repo.GetZhkMetro(ctx, int64(*zhk.MetroStationId))
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "zhkId": id, "err": err.Error()}).Warn("failed to get zhk metro")
		return domain.ZhkInfo{}, errors.New("Не удалось получить метро ЖК")
	}
	zhkMetro.Id = *zhk.MetroStationId

	// Собираем полную информацию о местоположении ЖК
	zhkAddress := domain.ZhkAddress{
		Address: zhk.Address,
		Metro:   zhkMetro,
	}

	// Получаем картинки ЖК
	zhkHeader, err := u.repo.GetZhkHeader(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkHeader failed")
		return domain.ZhkInfo{}, err
	}
	for i, img := range zhkHeader.Images {
		zhkHeader.Images[i] = "http://localhost:8003" + u.cfg.App.BaseImagesPath + img
	}

	// Собираем контактную информацию о ЖК
	zhkContacts := domain.ZhkContacts{Developer: zhk.Developer, Phone: zhk.Phone}

	// Получаем класс ЖК
	zhkCharacteristics, err := u.repo.GetZhkCharacteristics(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkCharacteristics failed")
		return domain.ZhkInfo{}, err
	}

	var zhkApartments domain.ZhkApartments

	// Подготавливаем данные о офферах
	prepareOfferDate(&zhkHeader, &zhkCharacteristics, &zhkApartments, zhkOffers)

	return domain.ZhkInfo{
		ID:              zhk.ID,
		Description:     zhk.Description,
		Address:         zhkAddress,
		Header:          zhkHeader,
		Contacts:        zhkContacts,
		Characteristics: zhkCharacteristics,
		Apartments:      zhkApartments,
	}, nil

}

func (u *zhkUsecase) GetAllZhk(ctx context.Context) ([]domain.ZhkInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	zhks, err := u.repo.GetAllZhk(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get AllZhk failed")
		return nil, err
	}

	var zhksInfo []domain.ZhkInfo
	for _, zhk := range zhks {
		zhkInfo, err := u.GetZhkInfo(ctx, zhk.ID)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(), "params": logger.LoggerFields{"id": zhk.ID}}).Error("GetZhk failed")
			return nil, err
		}
		zhksInfo = append(zhksInfo, zhkInfo)
	}

	return zhksInfo, nil
}

func prepareOfferDate(header *domain.ZhkHeader, characteristics *domain.ZhkCharacteristics, apartments *domain.ZhkApartments, offers *offerpb.GetOffersByZhkResponse) {
	var minPrice, maxPrice, minArea, maxArea, minFloors, maxFloors, minCeiling, maxCeiling int
	for _, offer := range offers.Offers {
		// Цена
		if offer.Price < int32(minPrice) {
			header.LowestPrice = int(offer.Price)
		}
		if offer.Price > int32(maxPrice) {
			header.HighestPrice = int(offer.Price)
		}
		// Площадь
		if offer.Area < int32(minArea) {
			characteristics.Square.LowestSquare = float64(offer.Area)
		}
		if offer.Area > int32(maxArea) {
			characteristics.Square.HighestSquare = float64(offer.Area)
		}
		// Этажность
		if offer.TotalFloors < int32(minFloors) {
			characteristics.Floors.LowestFloor = int(offer.TotalFloors)
		}
		if offer.TotalFloors > int32(maxFloors) {
			characteristics.Floors.HighestFloor = int(offer.TotalFloors)
		}
		// Высота потолков
		if offer.CeilingHeight < int32(minCeiling) {
			characteristics.CeilingHeight.LowestHeight = int(offer.CeilingHeight)
		}
		if offer.CeilingHeight > int32(maxCeiling) {
			characteristics.CeilingHeight.HighestHeight = int(offer.CeilingHeight)
		}
		characteristics.Decoration = append(characteristics.Decoration, int(offer.RenovationId))
	}
	// Оставляем уникальные виды отделок
	characteristics.Decoration = uniqueInts(characteristics.Decoration)

	aps := make(map[int]*domain.ZhkApartment)
	for _, offer := range offers.Offers {
		stats, ok := aps[int(offer.Rooms)]
		if !ok {
			aps[int(offer.Rooms)] = &domain.ZhkApartment{
				HighestPrice: int(offer.Price),
				LowestPrice: int(offer.Price),
				MinSquare: int(offer.Area),
				Offers: 1,
				Rooms: int(offer.Rooms),
			}
		} else {
			if offer.Price < int32(stats.LowestPrice) {
				stats.LowestPrice = int(offer.Price)
			}
			if offer.Price > int32(stats.HighestPrice) {
				stats.HighestPrice = int(offer.Price)
			}
			if offer.Area < int32(stats.MinSquare) {
				stats.MinSquare = int(offer.Area)
			}
			stats.Offers++
		}
	}
	apsResult := make([]domain.ZhkApartment, 0, len(aps))
	for _, apt := range aps {
		apsResult = append(apsResult, *apt)
	}
	apartments.Apartments = apsResult
}

func uniqueInts(input []int) []int {
    seen := make(map[int]struct{})
    result := make([]int, 0, len(input))
    for _, v := range input {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
    return result
}