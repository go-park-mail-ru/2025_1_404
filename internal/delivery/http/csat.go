package http

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	event := r.URL.Query().Get("event_name")
	if event == "" {
		utils.SendErrorResponse(w, "нет параметра event_name", http.StatusBadRequest)
		return
	}

	questions, err := h.UC.GetQuestionsByEvent(r.Context(), event)
	
	if err != nil {
		utils.SendErrorResponse(w, "event_name не найден", http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, questions, http.StatusOK)
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

	utils.SendJSONResponse(w, nil, http.StatusOK)
}

func (h *CsatHandler) GetAnswersByQuestion (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("question_id")
	if id == "" {
		utils.SendErrorResponse(w, "нет параметра question_id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, "не корректный question_id", http.StatusBadRequest)
		return
	}

	answers, err := h.UC.GetAnswersByQuestion(r.Context(), idInt)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении ответов", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, answers, http.StatusOK)
}

func (h *CsatHandler) GetEvents (w http.ResponseWriter, r *http.Request) {
	events, err := h.UC.GetEvents(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении events", http.StatusInternalServerError)
	}

	utils.SendJSONResponse(w, events, http.StatusOK)
}