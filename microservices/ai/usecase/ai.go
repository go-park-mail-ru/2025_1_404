package usecase

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/redis"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"google.golang.org/genai"
	"os"
	"time"
)

type aiUsecase struct {
	redisRepo redis.RedisRepo
	logger    logger.Logger
	cfg       *config.Config
}

type PropertiesType map[string]*genai.Schema

func NewAIUsecase(redisRepo redis.RedisRepo, logger logger.Logger, cfg *config.Config) *aiUsecase {
	return &aiUsecase{redisRepo: redisRepo, logger: logger, cfg: cfg}
}

func (u *aiUsecase) EvaluateOffer(ctx context.Context, offer domain.Offer) (*domain.EvaluationResult, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	offerSerialized, err := json.Marshal(offer)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: json marshal failed")
		return nil, err
	}

	evaluationResultText, err := u.redisRepo.Get(ctx, string(offerSerialized))
	if err != nil {
		SetProxy(u.cfg.Gemini.Proxy)
		client, err := GetClient(ctx, u.cfg.Gemini.Token)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: client creation failed")
			return nil, err
		}
		generateConfig := GetConfigForEvaluation(u.cfg.Gemini.EstimationPrompt)
		chat, err := client.Chats.Create(ctx, u.cfg.Gemini.Model, generateConfig, nil)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: chat create failed")
			return nil, err
		}
		result, err := chat.SendMessage(ctx, genai.Part{
			Text: string(offerSerialized),
		})
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: send message failed")
			return nil, err
		}
		evaluationResultText = result.Text()
		err = u.redisRepo.Set(ctx, string(offerSerialized), evaluationResultText, time.Hour*24)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: redis set failed")
			return nil, err
		}
	}
	var evaluationResult domain.EvaluationResult
	err = json.Unmarshal([]byte(evaluationResultText), &evaluationResult)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"requestID": requestID, "err": err.Error()}).Error("AI usecase: json unmarshal failed")
		return nil, err
	}
	return &evaluationResult, nil
}

func GetClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	geminiConfig := &genai.ClientConfig{
		APIKey: apiKey,
	}
	client, err := genai.NewClient(ctx, geminiConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetConfigForEvaluation(systemPrompt string) *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Role: "system",
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: PropertiesType{
				"market_price": {
					Type: genai.TypeObject,
					Properties: PropertiesType{
						"total": {
							Type:        genai.TypeNumber,
							Description: "The total market price",
						},
						"per_square_meter": {
							Type:        genai.TypeNumber,
							Description: "The market price per square meter",
						},
					},
				},
				"possible_cost_range": {
					Type: genai.TypeObject,
					Properties: PropertiesType{
						"min": {
							Type:        genai.TypeNumber,
							Description: "The minimum possible cost",
						},
						"max": {
							Type:        genai.TypeNumber,
							Description: "The maximum possible cost",
						},
					},
				},
			},
		},
	}
}

func SetProxy(proxy string) {
	os.Setenv("http_proxy", proxy)
	os.Setenv("https_proxy", proxy)
}
