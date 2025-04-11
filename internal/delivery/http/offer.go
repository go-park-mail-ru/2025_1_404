package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/pkg/content"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_404/domain"
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

const offerStatusDraft = 2
const offerStatusPublished = 1

func parseOfferFilter(r *http.Request) domain.OfferFilter {
	q := r.URL.Query()

	getInt := func(key string) *int {
		if val := q.Get(key); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil {
				return &parsed
			}
		}
		return nil
	}

	getString := func(key string) *string {
		if val := q.Get(key); val != "" {
			return &val
		}
		return nil
	}

	getBool := func(key string) *bool {
		if val := q.Get(key); val != "" {
			if val == "true" {
				b := true
				return &b
			} else if val == "false" {
				b := false
				return &b
			}
		}
		return nil
	}

	return domain.OfferFilter{
		MinArea:        getInt("min_area"),
		MaxArea:        getInt("max_area"),
		MinPrice:       getInt("min_price"),
		MaxPrice:       getInt("max_price"),
		Floor:          getInt("floor"),
		Rooms:          getInt("rooms"),
		Address:        getString("address"),
		RenovationID:   getInt("renovation_id"),
		PropertyTypeID: getInt("property_type_id"),
		PurchaseTypeID: getInt("purchase_type_id"),
		RentTypeID:     getInt("rent_type_id"),
		OfferTypeID:    getInt("offer_type_id"),
		NewBuilding:    getBool("new_building"),
		SellerID:       getInt("seller_id"),
	}
}

func hasFilter(f domain.OfferFilter) bool {
	return f.MinArea != nil || f.MaxArea != nil ||
		f.MinPrice != nil || f.MaxPrice != nil ||
		f.Floor != nil || f.Rooms != nil || f.Address != nil ||
		f.RenovationID != nil || f.PropertyTypeID != nil ||
		f.PurchaseTypeID != nil || f.RentTypeID != nil ||
		f.OfferTypeID != nil || f.NewBuilding != nil || f.SellerID != nil
}

func (h *OfferHandler) GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	filter := parseOfferFilter(r)

	// если хотя бы один фильтр задан — ищем по фильтру
	if hasFilter(filter) {
		offers, err := h.OfferUC.GetOffersByFilter(r.Context(), filter)
		if err != nil {
			utils.SendErrorResponse(w, "Ошибка при фильтрации объявлений", http.StatusInternalServerError)
			return
		}
		utils.SendJSONResponse(w, offers, http.StatusOK)
		return
	}

	// иначе — возвращаем все
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
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	var offer domain.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	offer.SellerID = userID
	offer.StatusID = offerStatusDraft

	id, err := h.OfferUC.CreateOffer(r.Context(), offer)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при создании", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]int{"id": id}, http.StatusCreated)
}

func (h *OfferHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	existingOffer, err := h.OfferUC.GetOfferByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, "Объявление не найдено", http.StatusNotFound)
		return
	}

	if existingOffer.SellerID != userID {
		utils.SendErrorResponse(w, "Нет доступа к обновлению этого объявления", http.StatusForbidden)
		return
	}

	var offer domain.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка в теле запроса", http.StatusBadRequest)
		return
	}

	offer.ID = id
	offer.SellerID = userID // Защита от подмены

	if err := h.OfferUC.UpdateOffer(r.Context(), offer); err != nil {
		utils.SendErrorResponse(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]string{"message": "Обновлено"}, http.StatusOK)
}

func (h *OfferHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	existingOffer, err := h.OfferUC.GetOfferByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, "Объявление не найдено", http.StatusNotFound)
		return
	}

	if existingOffer.SellerID != userID {
		utils.SendErrorResponse(w, "Нет доступа к удалению этого объявления", http.StatusForbidden)
		return
	}

	if err := h.OfferUC.DeleteOffer(r.Context(), id); err != nil {
		utils.SendErrorResponse(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]string{"message": "Удалено"}, http.StatusOK)
}

func (h *OfferHandler) UploadOfferImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	offerID, err := strconv.Atoi(vars["id"])
	if err != nil || offerID <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	offer, err := h.OfferUC.GetOfferByID(r.Context(), offerID)
	if err != nil || offer.SellerID != userID {
		utils.SendErrorResponse(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendErrorResponse(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(w, "Не удалось прочитать файл", http.StatusBadRequest)
		return
	}

	contentType, err := content.CheckImage(fileBytes)
	if err != nil {
		utils.SendErrorResponse(w, "Недопустимый формат изображения", http.StatusBadRequest)
		return
	}

	upload := filestorage.FileUpload{
		Name:        uuid.New().String() + "." + contentType,
		Size:        header.Size,
		File:        bytes.NewReader(fileBytes),
		ContentType: contentType,
	}

	imageID, err := h.OfferUC.SaveOfferImage(r.Context(), offerID, upload)
	if err != nil {
		utils.SendErrorResponse(w, "Ошибка при сохранении", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, map[string]any{
		"image_id": imageID,
	}, http.StatusCreated)
}

func (h *OfferHandler) PublishOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		utils.SendErrorResponse(w, "UserID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	offerID, err := strconv.Atoi(vars["id"])
	if err != nil || offerID <= 0 {
		utils.SendErrorResponse(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.OfferUC.PublishOffer(r.Context(), offerID, userID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, map[string]string{"message": "Объявление опубликовано"}, http.StatusOK)
}
