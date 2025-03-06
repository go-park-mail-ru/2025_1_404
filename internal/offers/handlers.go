package offers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_404/utils"
)

// GetOffersHandler Получение списка объявлений
func GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offers := GetOffers()
	utils.SendJSONResponse(w, offers, http.StatusOK)
}
