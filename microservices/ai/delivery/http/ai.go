package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai"
	"github.com/go-park-mail-ru/2025_1_404/microservices/ai/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"net/http"
)

type AIHandler struct {
	UC  ai.AIUsecase
	cfg *config.Config
}

func NewAIHandler(uc ai.AIUsecase, cfg *config.Config) *AIHandler {
	return &AIHandler{UC: uc, cfg: cfg}
}

func (h *AIHandler) EvaluateOffer(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	var offer domain.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	evaluationResult, err := h.UC.EvaluateOffer(r.Context(), offer)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при оценке предложения", http.StatusInternalServerError, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, evaluationResult, http.StatusOK, &h.cfg.App.CORS)
}
