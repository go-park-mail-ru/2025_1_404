package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type Offer struct {
	ID             int64
	SellerID       int64
	OfferTypeID    int
	MetroStationID *int
	RentTypeID     *int
	PurchaseTypeID *int
	PropertyTypeID int
	StatusID       int
	RenovationID   int
	ComplexID      *int
	Price          int
	Description    *string
	Floor          int
	TotalFloors    int
	Rooms          int
	Address        *string
	Flat           int
	Area           int
	CeilingHeight  int
	Longitude      string
	Latitude       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type offerRepository struct {
	db     database.DB
	logger logger.Logger
}

func NewOfferRepository(db database.DB, logger logger.Logger) *offerRepository {
	return &offerRepository{db: db, logger: logger}
}

const (
	createOfferSQL = `
		INSERT INTO kvartirum.Offer (
			seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		) RETURNING id;
	`

	getOfferByIDSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude, created_at, updated_at
		FROM kvartirum.Offer
		WHERE id = $1;
	`

	getOffersBySellerSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude, created_at, updated_at
		FROM kvartirum.Offer
		WHERE seller_id = $1;
	`

	getAllOffersSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude, created_at, updated_at
		FROM kvartirum.Offer;
	`

	getNotDraftOffersSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude, created_at, updated_at
		FROM kvartirum.Offer
		WHERE offer_status_id != 2;
	`

	updateOfferSQL = `
		UPDATE kvartirum.Offer
		SET offer_type_id = $1, metro_station_id = $2, rent_type_id = $3,
			purchase_type_id = $4, property_type_id = $5, offer_status_id = $6,
			renovation_id = $7, complex_id = $8, price = $9, description = $10,
			floor = $11, total_floors = $12, rooms = $13, address = $14,
			flat = $15, area = $16, ceiling_height = $17,  longitude = $18, latitude = $19
		WHERE id = $20;
	`

	deleteOfferSQL = `
		DELETE FROM kvartirum.Offer WHERE id = $1;
	`

	getOffersByZhkId = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, longitude, latitude, created_at, updated_at
		FROM kvartirum.Offer
		WHERE complex_id = $1;
	`

	getStations = `
		SELECT
			ms.id as station_id, ms.name as station_name, ml.color as color
		FROM kvartirum.MetroStation ms
		JOIN kvartirum.MetroLine ml ON ms.metro_line_id = ml.id;
	`
)

func (r *offerRepository) CreateOffer(ctx context.Context, o Offer) (int64, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	var id int64
	err := r.db.QueryRow(ctx, createOfferSQL,
		o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID,
		o.PurchaseTypeID, o.PropertyTypeID, o.StatusID, o.RenovationID,
		o.ComplexID, o.Price, o.Description, o.Floor, o.TotalFloors,
		o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight, o.Longitude, o.Latitude,
	).Scan(&id)

	logFields := logger.LoggerFields{"requestID": requestID, "query": createOfferSQL, "params": logger.LoggerFields{"seller_id": o.SellerID, "price": o.Price, "rooms": o.Rooms}, "success": err == nil}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query CreateOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query CreateOffer succeeded")
	}

	return id, err
}

func (r *offerRepository) GetOfferByID(ctx context.Context, id int64) (Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	var o Offer
	err := r.db.QueryRow(ctx, getOfferByIDSQL, id).Scan(
		&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
		&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
		&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
		&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
		&o.Longitude, &o.Latitude, &o.CreatedAt, &o.UpdatedAt,
	)

	logFields := logger.LoggerFields{"requestID": requestID, "query": getOfferByIDSQL, "params": logger.LoggerFields{"id": id}, "success": err == nil}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query GetOfferByID failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query GetOfferByID succeeded")
	}

	return o, err
}

func (r *offerRepository) GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	rows, err := r.db.Query(ctx, getOffersBySellerSQL, sellerID)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getOffersBySellerSQL, "params": logger.LoggerFields{"seller_id": sellerID}, "success": false, "err": err.Error()}).Error("SQL query GetOffersBySellerID failed")
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var o Offer
		err = rows.Scan(
			&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
			&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
			&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
			&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
			&o.Longitude, &o.Latitude, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getOffersBySellerSQL, "params": logger.LoggerFields{"seller_id": sellerID}, "success": true, "count": len(offers)}).Info("SQL query GetOffersBySellerID succeeded")

	return offers, nil
}

func (r *offerRepository) GetAllOffers(ctx context.Context) ([]Offer, error) {
	query := getNotDraftOffersSQL

	requestID := ctx.Value(utils.RequestIDKey)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": query, "success": false, "err": err.Error()}).Error("SQL query GetAllOffers failed")
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var o Offer
		err = rows.Scan(
			&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
			&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
			&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
			&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
			&o.Longitude, &o.Latitude, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": query, "success": true, "count": len(offers)}).Info("SQL query GetAllOffers succeeded")

	return offers, nil
}

func (r *offerRepository) GetOffersByFilter(ctx context.Context, f domain.OfferFilter, pUserId *int) ([]Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var (
		whereParts []string
		args       []any
		idx        = 1
	)

	addFilter := func(condition string, value any) {
		whereParts = append(whereParts, fmt.Sprintf(condition, idx))
		args = append(args, value)
		idx++
	}

	// Фильтры
	if f.MinArea != nil {
		addFilter("area >= $%d", *f.MinArea)
	}
	if f.MaxArea != nil {
		addFilter("area <= $%d", *f.MaxArea)
	}
	if f.MinPrice != nil {
		addFilter("price >= $%d", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		addFilter("price <= $%d", *f.MaxPrice)
	}
	if f.Floor != nil {
		addFilter("floor = $%d", *f.Floor)
	}
	if f.Rooms != nil {
		addFilter("rooms = $%d", *f.Rooms)
	}
	if f.Address != nil {
		addFilter("address ILIKE $%d", "%"+*f.Address+"%")
	}
	if f.RenovationID != nil {
		addFilter("renovation_id = $%d", *f.RenovationID)
	}
	if f.PropertyTypeID != nil {
		addFilter("property_type_id = $%d", *f.PropertyTypeID)
	}
	if f.PurchaseTypeID != nil {
		addFilter("purchase_type_id = $%d", *f.PurchaseTypeID)
	}
	if f.RentTypeID != nil {
		addFilter("rent_type_id = $%d", *f.RentTypeID)
	}
	if f.OfferTypeID != nil {
		addFilter("offer_type_id = $%d", *f.OfferTypeID)
	}
	if f.OnlyMe != nil && *f.OnlyMe && pUserId != nil {
		addFilter("seller_id = $%d", *pUserId)
	} else {
		addFilter("offer_status_id = $%d", 1)
		if f.SellerID != nil {
			addFilter("seller_id = $%d", *f.SellerID)
		}
	}
	if f.NewBuilding != nil {
		if *f.NewBuilding {
			whereParts = append(whereParts, "complex_id IS NOT NULL")
		} else {
			whereParts = append(whereParts, "complex_id IS NULL")
		}
	}

	query := strings.TrimRight(getAllOffersSQL, "\t\n;")

	if len(whereParts) > 0 {
		query += " WHERE " + strings.Join(whereParts, " AND ")
	}

	query += ";"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": query, "params": args, "success": false, "err": err.Error()}).Error("SQL query GetOffersByFilter failed")
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var o Offer
		err = rows.Scan(
			&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
			&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
			&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
			&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
			&o.Longitude, &o.Latitude, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": query, "params": args, "success": true, "count": len(offers)}).Info("SQL query GetOffersByFilter succeeded")

	return offers, nil
}

func (r *offerRepository) UpdateOffer(ctx context.Context, o Offer) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, updateOfferSQL,
		o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
		o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID,
		o.Price, o.Description, o.Floor, o.TotalFloors, o.Rooms,
		o.Address, o.Flat, o.Area, o.CeilingHeight,
		o.Longitude, o.Latitude, o.ID,
	)

	logFields := logger.LoggerFields{"requestID": requestID, "query": updateOfferSQL, "params": logger.LoggerFields{"id": o.ID, "price": o.Price}, "success": err == nil}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query UpdateOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query UpdateOffer succeeded")
	}

	return err
}

func (r *offerRepository) DeleteOffer(ctx context.Context, id int64) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, deleteOfferSQL, id)

	logFields := logger.LoggerFields{"requestID": requestID, "query": deleteOfferSQL, "params": logger.LoggerFields{"id": id}, "success": err == nil}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query DeleteOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query DeleteOffer succeeded")
	}

	return err
}

func (r *offerRepository) CreateImageAndBindToOffer(ctx context.Context, offerID int, uuid string) (int64, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var imageID int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO kvartirum.Image (uuid)
		VALUES ($1)
		RETURNING id;
	`, uuid).Scan(&imageID)

	if err != nil {
		r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "step": "insert image", "uuid": uuid, "err": err.Error()}).Error("Ошибка при вставке Image")
		return 0, err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO kvartirum.OfferImages (offer_id, image_id)
		VALUES ($1, $2);
	`, offerID, imageID)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "step": "bind to offer", "offer_id": offerID, "image_id": imageID, "err": err.Error()}).Error("Ошибка при вставке в OfferImages")
		return 0, err
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offerID, "image_id": imageID, "success": true}).Info("Изображение добавлено и связано с оффером")

	return imageID, nil
}

func (r *offerRepository) UpdateOfferStatus(ctx context.Context, offerID int, statusID int) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, `
		UPDATE kvartirum.Offer
		SET offer_status_id = $1
		WHERE id = $2;
	`, statusID, offerID)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offer_id": offerID, "status_id": statusID, "success": err == nil}).Info("SQL query UpdateOfferStatus")

	return err
}

func (r *offerRepository) GetOfferData(ctx context.Context, offer domain.Offer) (domain.OfferData, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var offerData domain.OfferData

	rows, err := r.db.Query(ctx, `
	SELECT
		i.id,
		i.uuid
	FROM kvartirum.OfferImages oi
	LEFT JOIN kvartirum.Image i ON oi.image_id = i.id
	WHERE oi.offer_id = $1
	ORDER BY oi.created_at;
	`, offer.ID)

	for rows.Next() {
		var offerImage domain.OfferImage
		err := rows.Scan(&offerImage.ID, &offerImage.Image)
		if err != nil {
			return domain.OfferData{}, err
		}
		offerData.Images = append(offerData.Images, offerImage)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offerID": offer.ID, "success": err == nil}).Info("SQL GetOfferImages")

	err = r.db.QueryRow(ctx, `
	SELECT
		ms.id as station_id, ms.name as station_name, ml.color as color
	FROM kvartirum.MetroStation ms
	LEFT JOIN kvartirum.MetroLine ml ON ms.metro_line_id = ml.id
	WHERE ms.id = $1;
	`, offer.MetroStationID).Scan(&offerData.Metro.Station, &offerData.Metro.Id, &offerData.Metro.Color)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offerID": offer.ID, "success": err == nil}).Info("SQL GetOfferStation")

	return offerData, nil
}

func (r *offerRepository) GetOfferImageWithUUID(ctx context.Context, imageID int64) (int64, string, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var offerID int64
	var uuid string

	err := r.db.QueryRow(ctx, `
		SELECT oi.offer_id, i.uuid
		FROM kvartirum.OfferImages oi
		JOIN kvartirum.Image i ON oi.image_id = i.id
		WHERE oi.image_id = $1;
	`, imageID).Scan(&offerID, &uuid)

	if err != nil {
		return 0, "", err
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "image_id": imageID, "offer_id": offerID, "uuid": uuid}).Info("Получена связь offer-image")

	return offerID, uuid, nil
}

func (r *offerRepository) DeleteOfferImage(ctx context.Context, imageID int64) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, `
		DELETE FROM kvartirum.OfferImages
		WHERE image_id = $1;
	`, imageID)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "image_id": imageID, "success": err == nil}).Info("SQL Delete OfferImage")

	return err
}

func (r *offerRepository) GetOffersByZhkId(ctx context.Context, zhkId int) ([]domain.Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var offers []domain.Offer

	rows, err := r.db.Query(ctx, getOffersByZhkId, zhkId)

	for rows.Next() {
		var o domain.Offer
		err := rows.Scan(
			&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
			&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
			&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
			&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
			&o.Longitude, &o.Latitude, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return offers, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "zhkID": zhkId, "success": err == nil}).Info("SQL GetOfferImages")

	return offers, nil
}

func (r *offerRepository) GetStations(ctx context.Context) ([]domain.Metro, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var metros []domain.Metro

	rows, err := r.db.Query(ctx, getStations)

	for rows.Next() {
		var metro domain.Metro
		err := rows.Scan(&metro.Id, &metro.Station, &metro.Color)
		if err != nil {
			return metros, err
		}
		metros = append(metros, metro)
	}

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "success": err == nil}).Info("SQL GetStations")

	return metros, nil
}
