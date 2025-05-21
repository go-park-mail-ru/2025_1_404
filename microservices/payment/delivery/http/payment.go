package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"net/http"
)

type PaymentHandler struct {
	UC  payment.PaymentUsecase
	cfg *config.Config
}

func NewPaymentHandler(uc payment.PaymentUsecase, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{UC: uc, cfg: cfg}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	var createRequest domain.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest, &h.cfg.App.CORS)
		return
	}

	evaluationResult, err := h.UC.CreatePayment(r.Context(), createRequest)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка создании ссылки на оплату", http.StatusInternalServerError, &h.cfg.App.CORS)
		return
	}

	utils.SendJSONResponse(w, evaluationResult, http.StatusOK, &h.cfg.App.CORS)
}
