package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/gorilla/mux"
)

type OfferHandler struct {
	OfferUC *usecase.OfferUsecase
}

func NewOfferHandler(uc *usecase.OfferUsecase) *OfferHandler {
	return &OfferHandler{OfferUC: uc}
}

// GetOffersHandler — получение списка объявлений
func (h *OfferHandler) GetOffersHandler(w http.ResponseWriter, r *http.Request) {

	offers, err := h.OfferUC.GetOffers(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении объявлений", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, offers, http.StatusOK)
}

func (h *OfferHandler) GetOfferByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	offer, err := h.OfferUC.GetOfferByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, "Объявление не найдено", http.StatusNotFound)
		return
	}

	utils.SendJSONResponse(w, offer, http.StatusOK)
}

func (h *OfferHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var offer domain.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	id, err := h.OfferUC.CreateOffer(r.Context(), offer)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]int{"id": id}, http.StatusCreated)
}

func (h *OfferHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var offer domain.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	offer.ID = id

	if err := h.OfferUC.UpdateOffer(r.Context(), offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]string{"message": "Обновлено"}, http.StatusOK)
}

func (h *OfferHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if err := h.OfferUC.DeleteOffer(r.Context(), id); err != nil {
		utils.SendErrorResponse(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]string{"message": "Удалено"}, http.StatusOK)
}
