package zhk

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_zhk_repo.go -package=mocks

type ZhkRepository interface {
	GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error)
	GetZhkHeader(ctx context.Context, zhk domain.Zhk) (domain.ZhkHeader, error)
	GetZhkCharacteristics(ctx context.Context, zhk domain.Zhk) (domain.ZhkCharacteristics, error)
	GetAllZhk(ctx context.Context) ([]domain.Zhk, error)
	GetZhkMetro(ctx context.Context, id int64) (domain.ZhkMetro, error)
}
