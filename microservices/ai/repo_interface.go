package ai

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/domain"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_ai_repo.go -package=mocks

type AIRepository interface {
	GetEvaluationOfOffer(ctx context.Context, offer domain.Offer) (domain.Offer, error)
}
