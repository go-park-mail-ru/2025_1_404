package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_zhk.go -package=mocks

type zhkUsecase interface {
	GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error)
	GetZhkInfo(ctx context.Context, zhk domain.Zhk) (domain.ZhkInfo, error)
	GetAllZhk(ctx context.Context) ([]domain.ZhkInfo, error)
}
