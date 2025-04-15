package http

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/gorilla/mux"
)

type ZhkHandler struct {
	UC zhkUsecase
}

func NewZhkHandler(uc zhkUsecase) *ZhkHandler {
	return &ZhkHandler{UC: uc}
}

func (h *ZhkHandler) GetZhkInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, "неккоректный id ЖК", http.StatusBadRequest)
		return
	}

	zhk, err := h.UC.GetZhkByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, "ЖК с таким id не найден", http.StatusNotFound)
		return
	}

	zhkInfo, err := h.UC.GetZhkInfo(r.Context(), zhk)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении ЖК", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, zhkInfo, http.StatusOK)
}

func (h *ZhkHandler) GetAllZhk(w http.ResponseWriter, r *http.Request) {
	zhks, err := h.UC.GetAllZhk(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при получении списка ЖК", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, zhks, http.StatusOK)
}
