package http

import (
	"github.com/go-park-mail-ru/2025_1_404/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"net/http"
)

// GetOffersHandler Получение списка объявлений
func GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offers := usecase.GetOffers()
	utils.SendJSONResponse(w, offers, http.StatusOK)
}
