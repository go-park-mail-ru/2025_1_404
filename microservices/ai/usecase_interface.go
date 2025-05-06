package ai

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_ai.go -package=mocks

type AIUsecase interface {
	EvaluateOffer(ctx context.Context, offer domain.Offer) (*domain.EvaluationResult, error)
}
