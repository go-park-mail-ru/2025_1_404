package repository

import (
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"golang.org/x/net/context"
)

type csatRepository struct {
	db     repository.DB
	logger logger.Logger
}

func NewCsatRepository(db repository.DB, logger logger.Logger) *csatRepository {
	return &csatRepository{db: db, logger: logger}
}

const (
	getQuestionsByEvent = `
		SELECT id, text, csat_id
		FROM csat.Question q
		WHERE q.csat_id = (SELECT id from csat.csat WHERE event = $1);
	`

	createAnswer = `
		INSERT INTO csat.Answer (rating, question_id) VALUES
			($1, $2)
		RETURNING id;
	`

	getAnswerStat = `
		SELECT
			COUNT(*) AS total_answers,
			COALESCE(AVG(rating), 0) AS avg_rating,
			COUNT(*) FILTER (WHERE rating = 1) AS one_star_amount,
			COALESCE(COUNT(*) FILTER (WHERE rating = 1) * 100.0 / NULLIF(COUNT(*), 0), 0) AS one_star_percentage,
			COUNT(*) FILTER (WHERE rating = 2) AS two_star_amount,
			COALESCE(COUNT(*) FILTER (WHERE rating = 2) * 100.0 / NULLIF(COUNT(*), 0), 0) AS two_star_percentage,
			COUNT(*) FILTER (WHERE rating = 3) AS three_star_amount,
			COALESCE(COUNT(*) FILTER (WHERE rating = 3) * 100.0 / NULLIF(COUNT(*), 0), 0) AS three_star_percentage,
			COUNT(*) FILTER (WHERE rating = 4) AS four_star_amount,
			COALESCE(COUNT(*) FILTER (WHERE rating = 4) * 100.0 / NULLIF(COUNT(*), 0), 0) AS four_star_percentage,
			COUNT(*) FILTER (WHERE rating = 5) AS five_star_amount,
			COALESCE(COUNT(*) FILTER (WHERE rating = 5) * 100.0 / NULLIF(COUNT(*), 0), 0) AS five_star_percentage
		FROM csat.Answer a
		WHERE a.question_id = $1;
	`
)

func (repo *csatRepository) GetQuestionsByEvent(ctx context.Context, event string) ([]domain.QuestionDTO, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rows, err := repo.db.Query(ctx, getQuestionsByEvent, event)
	if err != nil {
		repo.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getQuestionsByEvent, "params": logger.LoggerFields{"event": event}, "err": err.Error()}).Error("SQL Query: get questionsByEvent failed")
		return []domain.QuestionDTO{}, err
	}
	defer rows.Close()

	var questions []domain.QuestionDTO
	for rows.Next() {
		var question domain.QuestionDTO
		err := rows.Scan(&question.ID, &question.Text, &question.CsatID)
		if err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (repo *csatRepository) AddAnswerToQuestion(ctx context.Context, answer domain.AnswerDTO) error {
	requestID := ctx.Value(utils.RequestIDKey)

	var id int64
	err := repo.db.QueryRow(ctx, createAnswer, answer.Rating, answer.QuestionID).Scan(&id)

	repo.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"params": logger.LoggerFields{
			"question_id": answer.QuestionID,
			"rating":      answer.Rating,
		},
		"success": err == nil,
	}).Info("SQL Query add answer to question")

	return err
}

func (repo *csatRepository) GetAnswersByQuestion(ctx context.Context, questionID int64) (domain.AnswersStat, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var answers domain.AnswersStat
	err := repo.db.QueryRow(ctx, getAnswerStat, questionID).Scan(
		&answers.TotalAnswers, &answers.AvgRating, &answers.OneStarStat.Amount, &answers.OneStarStat.Percentage,
		&answers.TwoStarStat.Amount, &answers.TwoStarStat.Percentage, &answers.ThreeStarStat.Amount, &answers.ThreeStarStat.Percentage,
		&answers.FourStarStat.Amount, &answers.FourStarStat.Percentage, &answers.FiveStarStat.Amount, &answers.FiveStarStat.Percentage,
	)

	repo.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getAnswerStat, "params": logger.LoggerFields{"questionID": questionID}, "success": err == nil}).Info("SQL Query: get answers stat")

	return answers, err
}
