package usecase

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type zhkUsecase struct {
	repo   zhkRepository
	logger logger.Logger
}

func NewZhkUsecase(repo zhkRepository, logger logger.Logger) *zhkUsecase {
	return &zhkUsecase{repo: repo, logger: logger}
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

func (u *zhkUsecase) GetZhkInfo(ctx context.Context, zhk domain.Zhk) (domain.ZhkInfo, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	log.Println(zhk)

	zhkHeader, err := u.repo.GetZhkHeader(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkHeader failed")
		return domain.ZhkInfo{}, err
	}
	for i, img := range zhkHeader.Images {
		zhkHeader.Images[i] = utils.BasePath + utils.ImagesPath + img
	}

	zhkContacts := domain.ZhkContacts{Developer: zhk.Developer, Phone: zhk.Phone}

	zhkCharacteristics, err := u.repo.GetZhkCharacteristics(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkCharacteristics failed")
		return domain.ZhkInfo{}, err
	}

	zhkApartments, err := u.repo.GetZhkApartments(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkApartments failed")
		return domain.ZhkInfo{}, err
	}

	zhkReviews, err := u.repo.GetZhkReviews(ctx, zhk)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("Get ZhkReviews failed")
		return domain.ZhkInfo{}, err
	}
	for i, review := range zhkReviews.Reviews {
		if review.Avatar != "" {
			zhkReviews.Reviews[i].Avatar = utils.BasePath + utils.ImagesPath + review.Avatar
		}
	}

	return domain.ZhkInfo{
		ID:              zhk.ID,
		Description:     zhk.Description,
		Address:         zhk.Address,
		Header:          zhkHeader,
		Contacts:        zhkContacts,
		Characteristics: zhkCharacteristics,
		Apartments:      zhkApartments,
		Reviews:         zhkReviews,
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
		zhkInfo, err := u.GetZhkInfo(ctx, zhk)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error(),
				"params": logger.LoggerFields{"id": zhk.ID}}).Error("GetZhk failed")
			return nil, err
		}
		zhksInfo = append(zhksInfo, zhkInfo)
	}

	return zhksInfo, nil
}
