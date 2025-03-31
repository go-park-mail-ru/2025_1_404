package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type OfferHandler struct {
	OfferUC *usecase.OfferUsecase
}

func NewOfferHandler(uc *usecase.OfferUsecase) *OfferHandler {
	return &OfferHandler{OfferUC: uc}
}

// GetOffersHandler — получение списка объявлений
func (h *OfferHandler) GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offers, err := h.OfferUC.GetOffers(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении объявлений", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, offers, http.StatusOK)
}
