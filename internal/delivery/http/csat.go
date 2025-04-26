package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type CsatHandler struct {
	UC csatUsecase
}

func NewCsatHandler (uc csatUsecase) *CsatHandler {
	return &CsatHandler{UC: uc}
}

func (h *CsatHandler) GetQuestionsByEvent (w http.ResponseWriter, r *http.Request) {
	event := r.URL.Query().Get("event")
	if event == "" {
		utils.SendErrorResponse(w, "нет параметра event", http.StatusBadRequest)
		return
	}

	questions, err := h.UC.GetQuestionsByEvent(r.Context(), event)
	
	if err != nil {
		utils.SendErrorResponse(w, "event не найден", http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, questions, http.StatusFound)
}

func (h *CsatHandler) AddAnswerToQuestion (w http.ResponseWriter, r *http.Request) {
	var req domain.AnswerDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	err := h.UC.AddAnswerToQuestion(r.Context(), req)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании ответа", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, nil, http.StatusCreated)
}

func (h *CsatHandler) GetAnswersByQuestion (w http.ResponseWriter, r *http.Request) {
	var req domain.QuestionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	answers, err := h.UC.GetAnswersByQuestion(r.Context(), req.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении ответов", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, answers, http.StatusOK)
}
