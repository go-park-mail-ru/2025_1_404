package zhk

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_zhk.go -package=mocks

type ZhkUsecase interface {
	GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error)
	GetZhkInfo(ctx context.Context, id int64) (domain.ZhkInfo, error)
	GetAllZhk(ctx context.Context) ([]domain.ZhkInfo, error)
}
