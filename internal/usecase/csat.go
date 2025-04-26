package usecase

import (
	"context"
	"math"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type csatUsecase struct {
	repo   csatRepo
	logger logger.Logger
}

func NewCsatUsecase(repo csatRepo, logger logger.Logger) *csatUsecase {
	return &csatUsecase{repo: repo, logger: logger}
}

func (u *csatUsecase) GetQuestionsByEvent(ctx context.Context, event string) ([]domain.QuestionDTO, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	questions, err := u.repo.GetQuestionsByEvent(ctx, event)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("CSAT usecase: get questions by event failed")
		return []domain.QuestionDTO{}, err
	}

	return questions, nil
}

func (u *csatUsecase) AddAnswerToQuestion(ctx context.Context, answer domain.AnswerDTO) error {
	requestID := ctx.Value(utils.RequestIDKey)

	err := u.repo.AddAnswerToQuestion(ctx, answer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("CSAT usecase: add answer failed")
		return err
	}

	return nil
}

func (u *csatUsecase) GetAnswersByQuestion(ctx context.Context, questionID int64) (domain.AnswersStat, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	answers, err := u.repo.GetAnswersByQuestion(ctx, questionID)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("CSAT usecase: get answers failed")
		return domain.AnswersStat{}, err
	}

	digitsAfterPoint := 0.1
	answers.AvgRating = truncateNaive(answers.AvgRating, digitsAfterPoint)
	answers.OneStarStat.Percentage = truncateNaive(answers.OneStarStat.Percentage, digitsAfterPoint)
	answers.TwoStarStat.Percentage = truncateNaive(answers.TwoStarStat.Percentage, digitsAfterPoint)
	answers.ThreeStarStat.Percentage = truncateNaive(answers.ThreeStarStat.Percentage, digitsAfterPoint)
	answers.FourStarStat.Percentage = truncateNaive(answers.FourStarStat.Percentage, digitsAfterPoint)
	answers.FiveStarStat.Percentage = truncateNaive(answers.FiveStarStat.Percentage, digitsAfterPoint)

	return answers, nil
}

func (u *csatUsecase) GetEvents(ctx context.Context) (domain.EventList, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	events, err := u.repo.GetEvents(ctx)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("CSAT usecase: get events failed")
		return domain.EventList{}, err
	}

	return events, nil
}

func truncateNaive(f float64, unit float64) float64 {
	return math.Trunc(f/unit) * unit
}
