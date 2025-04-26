package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/domain"
)

type csatUsecase interface {
	GetQuestionsByEvent (ctx context.Context, event string) ([]domain.QuestionDTO, error)
	AddAnswerToQuestion (ctx context.Context, answer domain.AnswerDTO) (error)
	GetAnswersByQuestion (ctx context.Context, questionID int64) (domain.AnswersStat, error)
	GetEvents(ctx context.Context) (domain.EventList, error)
}