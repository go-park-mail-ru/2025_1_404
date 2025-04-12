// internal/repository/repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -source repository.go -destination=mocks/mock_repository.go -package=mocks

type User struct {
	ID                 int64
	ImageID            *int64
	FirstName          string
	LastName           string
	Email              string
	Password           string
	LastNotificationID *int
	TokenVersion       int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

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
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Repository interface {
	// --- Users ---
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateImage(ctx context.Context, file filestorage.FileUpload) error
	GetImageByID(ctx context.Context, id sql.NullInt64) (string, error)
	DeleteUserImage(ctx context.Context, id int64) error

	// --- Offers ---
	CreateOffer(ctx context.Context, offer Offer) (int64, error)
	GetOfferByID(ctx context.Context, id int64) (Offer, error)
	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error)
	GetAllOffers(ctx context.Context) ([]Offer, error)
	GetOffersByFilter(ctx context.Context, f domain.OfferFilter) ([]Offer, error)
	UpdateOffer(ctx context.Context, offer Offer) error
	DeleteOffer(ctx context.Context, id int64) error
 	CreateImageAndBindToOffer(ctx context.Context, offerID int, uuid string) (int64, error)
	UpdateOfferStatus(ctx context.Context, offerID int, statusID int) error
	GetOfferData(ctx context.Context, offer domain.Offer) (domain.OfferData, error)
	GetOfferImageWithUUID(ctx context.Context, imageID int64) (int64, string, error)
	DeleteOfferImage(ctx context.Context, imageID int64) error

	// --- Zhk ---
	GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error)
	GetZhkHeader(ctx context.Context, zhk domain.Zhk) (domain.ZhkHeader, error)
	GetZhkCharacteristics(ctx context.Context, zhk domain.Zhk) (domain.ZhkCharacteristics, error)
	GetZhkApartments(ctx context.Context, zhk domain.Zhk) (domain.ZhkApartments, error)
	GetZhkReviews(ctx context.Context, zhk domain.Zhk) (domain.ZhkReviews, error)
	GetAllZhk(ctx context.Context) ([]domain.Zhk, error)
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

type repository struct {
	db     DB
	logger logger.Logger
}

func NewRepository(db DB, logger logger.Logger) Repository {
	return &repository{db: db, logger: logger}
}

// region --- USERS ---

const (
	createUserSQL = `
		INSERT INTO kvartirum.Users (
			image_id, first_name, last_name, email, password, token_version
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	getUserByEmailSQL = `
		SELECT
			u.id, 
			COALESCE(i.uuid, '') as image,
			u.first_name, u.last_name, u.email, u.password
		FROM kvartirum.Users u
		LEFT JOIN kvartirum.Image i on u.image_id = i.id
		WHERE u.email = $1;
	`

	getUserByIDSQL = `
		SELECT
			u.id, 
			COALESCE(i.uuid, '') as image, 
			u.first_name, u.last_name, u.email, u.password
		FROM kvartirum.Users u
		LEFT JOIN kvartirum.Image i on u.image_id = i.id
		WHERE u.id = $1;
	`

	updateUserSQL = `
		UPDATE kvartirum.Users
		SET
			image_id = (SELECT id FROM kvartirum.Image WHERE uuid = $1),
			first_name = $2, last_name = $3, email = $4
		WHERE id = $5
		RETURNING 
			id, first_name, last_name, email,
		 	COALESCE ((SELECT uuid FROM kvartirum.Image WHERE id = image_id), '') as image_uuid;
	`

	deleteUserSQL = `
		DELETE FROM kvartirum.Users WHERE id = $1;
	`

	createImageSQL = `
		INSERT INTO kvartirum.Image (uuid) VALUES ($1);
	`
	getImageByUUIDSQL = `
		SELECT id FROM kvartirum.Image WHERE uuid = $1;
	`

	getImageByIDSQL = `
		SELECT uuid from kvartirum.Image WHERE id = $1;
	`

	deleteUserImageSQL = `
		DELETE FROM kvartirum.Image where id = (
			SELECT image_id from kvartirum.Users where id = $1
		);
	`
)

func (r *repository) CreateUser(ctx context.Context, u User) (int64, error) {
	var id int64
	requestID := ctx.Value(utils.RequestIDKey)
	err := r.db.QueryRow(ctx, createUserSQL,
		u.ImageID, u.FirstName, u.LastName, u.Email, u.Password, u.TokenVersion,
	).Scan(&id)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID, "query": createUserSQL,
		"params": logger.LoggerFields{
			"name":          u.FirstName,
			"last_name":     u.LastName,
			"email":         u.Email,
			"token_version": u.TokenVersion,
			"image_id":      u.ImageID,
		}, "success": err == nil}).Info("SQL query CreateUser")

	return id, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var u domain.User
	err := r.db.QueryRow(ctx, getUserByEmailSQL, email).Scan(
		&u.ID, &u.Image, &u.FirstName, &u.LastName, &u.Email, &u.Password,
	)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getUserByEmailSQL,
		"params": logger.LoggerFields{"email": email}, "success": err == nil}).Info("SQL query GetUserByEmail")

	return u, err
}

func (r *repository) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var u domain.User
	err := r.db.QueryRow(ctx, getUserByIDSQL, id).Scan(
		&u.ID, &u.Image, &u.FirstName, &u.LastName, &u.Email, &u.Password,
	)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     getUserByEmailSQL,
		"params": logger.LoggerFields{
			"id": id,
		},
		"success": err == nil,
	}).Info("SQL query GetUserByID")

	return u, err
}

func (r *repository) UpdateUser(ctx context.Context, u domain.User) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var updatedUser domain.User
	err := r.db.QueryRow(ctx, updateUserSQL,
		u.Image, u.FirstName, u.LastName, u.Email, u.ID,
	).Scan(&updatedUser.ID, &updatedUser.FirstName, &updatedUser.LastName,
		&updatedUser.Email, &updatedUser.Image)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID, "query": updateUserSQL,
		"params": logger.LoggerFields{
			"name":       u.FirstName,
			"last_name":  u.LastName,
			"email":      u.Email,
			"image_path": u.Image,
		}, "success": err == nil}).Info("SQL query UpdateUser")

	return updatedUser, err
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, deleteUserSQL, id)
	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     deleteUserSQL,
		"params": logger.LoggerFields{
			"id": id,
		},
		"success": err == nil,
	})
	return err
}

func (r *repository) CreateImage(ctx context.Context, file filestorage.FileUpload) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, createImageSQL, file.Name)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     createImageSQL,
		"params": logger.LoggerFields{
			"name": file.Name,
		},
		"success": err == nil,
	})

	return err
}

func (r *repository) GetImageByUUID(ctx context.Context, uuid string) (sql.NullInt64, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var id sql.NullInt64
	err := r.db.QueryRow(ctx, getImageByUUIDSQL, uuid).Scan(&id)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     getImageByUUIDSQL,
		"params": logger.LoggerFields{
			"name": uuid,
		},
		"success": err == nil,
	})

	return id, err
}

func (r *repository) GetImageByID(ctx context.Context, id sql.NullInt64) (string, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var fileName string
	if !id.Valid {
		return "", nil
	}
	err := r.db.QueryRow(ctx, getImageByIDSQL, id.Int64).Scan(&fileName)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID, "query": getImageByIDSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("SQL GetIMageByID")

	return fileName, err

}

func (r *repository) DeleteUserImage(ctx context.Context, id int64) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, deleteUserImageSQL, id)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": deleteUserSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("SQL DeleteImage")

	return err
}

// endregion

// region --- OFFERS ---

const (
	createOfferSQL = `
		INSERT INTO kvartirum.Offer (
			seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id;
	`

	getOfferByIDSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, created_at, updated_at
		FROM kvartirum.Offer
		WHERE id = $1;
	`

	getOffersBySellerSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, created_at, updated_at
		FROM kvartirum.Offer
		WHERE seller_id = $1;
	`

	getAllOffersSQL = `
		SELECT id, seller_id, offer_type_id, metro_station_id, rent_type_id,
			purchase_type_id, property_type_id, offer_status_id, renovation_id,
			complex_id, price, description, floor, total_floors, rooms,
			address, flat, area, ceiling_height, created_at, updated_at
		FROM kvartirum.Offer;
	`

	updateOfferSQL = `
		UPDATE kvartirum.Offer
		SET offer_type_id = $1, metro_station_id = $2, rent_type_id = $3,
			purchase_type_id = $4, property_type_id = $5, offer_status_id = $6,
			renovation_id = $7, complex_id = $8, price = $9, description = $10,
			floor = $11, total_floors = $12, rooms = $13, address = $14,
			flat = $15, area = $16, ceiling_height = $17
		WHERE id = $18;
	`

	deleteOfferSQL = `
		DELETE FROM kvartirum.Offer WHERE id = $1;
	`
)

func (r *repository) CreateOffer(ctx context.Context, o Offer) (int64, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	var id int64
	err := r.db.QueryRow(ctx, createOfferSQL,
		o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID,
		o.PurchaseTypeID, o.PropertyTypeID, o.StatusID, o.RenovationID,
		o.ComplexID, o.Price, o.Description, o.Floor, o.TotalFloors,
		o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
	).Scan(&id)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     createOfferSQL,
		"params": logger.LoggerFields{
			"seller_id": o.SellerID,
			"price":     o.Price,
			"rooms":     o.Rooms,
		},
		"success": err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query CreateOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query CreateOffer succeeded")
	}

	return id, err
}

func (r *repository) GetOfferByID(ctx context.Context, id int64) (Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	var o Offer
	err := r.db.QueryRow(ctx, getOfferByIDSQL, id).Scan(
		&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
		&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
		&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
		&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
		&o.CreatedAt, &o.UpdatedAt,
	)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getOfferByIDSQL,
		"params":    logger.LoggerFields{"id": id},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query GetOfferByID failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query GetOfferByID succeeded")
	}

	return o, err
}

func (r *repository) GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	rows, err := r.db.Query(ctx, getOffersBySellerSQL, sellerID)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"query":     getOffersBySellerSQL,
			"params":    logger.LoggerFields{"seller_id": sellerID},
			"success":   false,
			"err":       err.Error(),
		}).Error("SQL query GetOffersBySellerID failed")
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
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     getOffersBySellerSQL,
		"params":    logger.LoggerFields{"seller_id": sellerID},
		"success":   true,
		"count":     len(offers),
	}).Info("SQL query GetOffersBySellerID succeeded")

	return offers, nil
}

func (r *repository) GetAllOffers(ctx context.Context) ([]Offer, error) {
	requestID := ctx.Value(utils.RequestIDKey)
	rows, err := r.db.Query(ctx, getAllOffersSQL)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"query":     getAllOffersSQL,
			"success":   false,
			"err":       err.Error(),
		}).Error("SQL query GetAllOffers failed")
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
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     getAllOffersSQL,
		"success":   true,
		"count":     len(offers),
	}).Info("SQL query GetAllOffers succeeded")

	return offers, nil
}

func (r *repository) GetOffersByFilter(ctx context.Context, f domain.OfferFilter) ([]Offer, error) {
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
	if f.SellerID != nil {
		addFilter("seller_id = $%d", *f.SellerID)
	}
	if f.NewBuilding != nil {
		if *f.NewBuilding {
			whereParts = append(whereParts, "complex_id IS NOT NULL")
		} else {
			whereParts = append(whereParts, "complex_id IS NULL")
		}
	}

	query := getAllOffersSQL

	if len(whereParts) > 0 {
		query += " WHERE " + strings.Join(whereParts, " AND ")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"query":     query,
			"params":    args,
			"success":   false,
			"err":       err.Error(),
		}).Error("SQL query GetOffersByFilter failed")
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
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     query,
		"params":    args,
		"success":   true,
		"count":     len(offers),
	}).Info("SQL query GetOffersByFilter succeeded")

	return offers, nil
}

func (r *repository) UpdateOffer(ctx context.Context, o Offer) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, updateOfferSQL,
		o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
		o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID,
		o.Price, o.Description, o.Floor, o.TotalFloors, o.Rooms,
		o.Address, o.Flat, o.Area, o.CeilingHeight, o.ID,
	)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     updateOfferSQL,
		"params": logger.LoggerFields{
			"id":    o.ID,
			"price": o.Price,
		},
		"success": err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query UpdateOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query UpdateOffer succeeded")
	}

	return err
}

func (r *repository) DeleteOffer(ctx context.Context, id int64) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, deleteOfferSQL, id)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     deleteOfferSQL,
		"params":    logger.LoggerFields{"id": id},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("SQL query DeleteOffer failed")
	} else {
		r.logger.WithFields(logFields).Info("SQL query DeleteOffer succeeded")
	}

	return err
}

func (r *repository) CreateImageAndBindToOffer(ctx context.Context, offerID int, uuid string) (int64, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var imageID int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO kvartirum.Image (uuid)
		VALUES ($1)
		RETURNING id;
	`, uuid).Scan(&imageID)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"step":      "insert image",
			"uuid":      uuid,
			"err":       err.Error(),
		}).Error("Ошибка при вставке Image")
		return 0, err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO kvartirum.OfferImages (offer_id, image_id)
		VALUES ($1, $2);
	`, offerID, imageID)
	if err != nil {
		r.logger.WithFields(logger.LoggerFields{
			"requestID": requestID,
			"step":      "bind to offer",
			"offer_id":  offerID,
			"image_id":  imageID,
			"err":       err.Error(),
		}).Error("Ошибка при вставке в OfferImages")
		return 0, err
	}

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"offer_id":  offerID,
		"image_id":  imageID,
		"success":   true,
	}).Info("Изображение добавлено и связано с оффером")

	return imageID, nil
}

func (r *repository) UpdateOfferStatus(ctx context.Context, offerID int, statusID int) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, `
		UPDATE kvartirum.Offer
		SET offer_status_id = $1
		WHERE id = $2;
	`, statusID, offerID)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"offer_id":  offerID,
		"status_id": statusID,
		"success":   err == nil,
	}).Info("SQL query UpdateOfferStatus")

	return err
}

func (r *repository) GetOfferData(ctx context.Context, offer domain.Offer) (domain.OfferData, error) {
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

	err = r.db.QueryRow(ctx, getUserByIDSQL, offer.SellerID).Scan(
		new(int64), &offerData.Seller.Avatar, &offerData.Seller.FirstName,
		&offerData.Seller.LastName, new(string), new(string))

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offerID": offer.ID, "success": err == nil}).Info("SQL GetOfferSeller")

	err = r.db.QueryRow(ctx, `
	SELECT
		ms.name as station_name,
		ml.name as line_name
	FROM kvartirum.MetroStation ms
	LEFT JOIN kvartirum.MetroLine ml ON ms.metro_line_id = ml.id
	WHERE ms.id = $1;
	`,offer.MetroStationID).Scan(&offerData.Metro.Station, &offerData.Metro.Line)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "offerID": offer.ID, "success": err == nil}).Info("SQL GetOfferStation")

	return offerData, nil
}

func (r *repository) GetOfferImageWithUUID(ctx context.Context, imageID int64) (int64, string, error) {
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

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"image_id":  imageID,
		"offer_id":  offerID,
		"uuid":      uuid,
	}).Info("Получена связь offer-image")

	return offerID, uuid, nil
}

func (r *repository) DeleteOfferImage(ctx context.Context, imageID int64) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, `
		DELETE FROM kvartirum.OfferImages
		WHERE image_id = $1;
	`, imageID)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"image_id":  imageID,
		"success":   err == nil,
	}).Info("SQL Delete OfferImage")

	return err
}

// endregion

// region --- ZHK ---

const (
	getZhkByIDSQL = `
	SELECT id, class_id, name, developer, phone_number, address, description
	FROM kvartirum.HousingComplex
	WHERE id = $1;
	`
	getZhkHeaderSQL = `
	SELECT
		COALESCE (MIN(o.price), 0) as lowest_price,
		COALESCE (MAX(o.price), 0) as highest_price,
		COALESCE(ARRAY_AGG(DISTINCT img.uuid) FILTER (WHERE img.uuid IS NOT NULL), '{}') AS images,
		COUNT (DISTINCT img.id) as images_size
	FROM kvartirum.housingcomplex hc
	LEFT JOIN kvartirum.Offer o ON o.complex_id = hc.id
	LEFT JOIN kvartirum.HousingComplexImages hci on hci.housing_complex_id = hc.id
	LEFT JOIN kvartirum.Image img on img.id = hci.image_id
	WHERE hc.id = $1
	GROUP BY hc.id, hc.name;
	`

	getZhkCharacteristicsSQL = `
	SELECT
		hcc.name as class_name,
		COALESCE (ARRAY_AGG(DISTINCT offerrenovation.name) 
			FILTER (WHERE offerrenovation.name IS NOT NULL), '{}') AS decoration,
		COALESCE (MAX(o.ceiling_height), 0) as max_ceiling_height,
		COALESCE (MIN(o.ceiling_height), 0) as min_ceiling_height,
		COALESCE(MAX(o.total_floors), 0) AS max_floors,
		COALESCE(MIN(o.total_floors), 0) AS min_floors,
		COALESCE(MAX(o.area), 0) AS max_area,
		COALESCE(MIN(o.area), 0) AS min_area
		FROM kvartirum.housingcomplex hc
		LEFT JOIN kvartirum.housingcomplexclass hcc ON hcc.id = hc.class_id
		LEFT JOIN kvartirum.offer o ON o.complex_id = hc.id
		LEFT JOIN kvartirum.offerrenovation ON offerrenovation.id = o.renovation_id
		WHERE hc.id = $1
		GROUP BY hcc.name;
			
	`

	getZhkApartmentsSQL = `
	SELECT
		o.rooms as rooms,
		COALESCE (MIN(o.price), 0) as lowest_price,
		COALESCE (MAX(o.price), 0) as highest_price,
		COALESCE (MIN(o.area), 0) as min_square,
		COUNT (*) as offers
	FROM kvartirum.offer o
	WHERE o.complex_id = $1
	GROUP BY o.rooms
	ORDER BY o.rooms;
	`

	getZhkReviewsSQL = `
	SELECT
		COALESCE(img.uuid, '')  as avatar,
		u.first_name as first_name,
		u.last_name as last_name,
		r.rating,
		r.comment
	FROM kvartirum.housingcomplexreview r
	JOIN kvartirum.users u ON u.id = r.user_id
	LEFT JOIN kvartirum.image img ON img.id = u.image_id
	WHERE r.housing_complex_id = $1
	ORDER BY r.created_at DESC
	`

	getZhkReviewsParamsSQL = `
	SELECT 
		COUNT(*) AS quantity,
		COALESCE(AVG(r.rating), 0) AS avg_rating
	FROM kvartirum.housingcomplexreview r
	WHERE r.housing_complex_id = $1;
	`

	getAllZhkSQL = `
	SELECT 
		id, class_id, name, developer, phone_number, address, description
	FROM kvartirum.housingcomplex;
	`
)

func (r *repository) GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var zhk domain.Zhk
	err := r.db.QueryRow(ctx, getZhkByIDSQL, id).Scan(
		&zhk.ID, &zhk.ClassID, &zhk.Name, &zhk.Developer,
		&zhk.Phone, &zhk.Address, &zhk.Description,
	)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getZhkByIDSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("GetZhkByID")

	fmt.Println("REPO", zhk)

	return zhk, err
}

func (r *repository) GetZhkHeader(ctx context.Context, zhk domain.Zhk) (domain.ZhkHeader, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	header := domain.ZhkHeader{Name: zhk.Name}

	err := r.db.QueryRow(ctx, getZhkHeaderSQL, zhk.ID).Scan(
		&header.LowestPrice, &header.HighestPrice, &header.Images, &header.ImagesSize,
	)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getZhkHeaderSQL,
		"params":    logger.LoggerFields{"id": zhk.ID},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetZhkHeader failed")
	} else {
		r.logger.WithFields(logFields).Info("GetZhkHeader succeeded")
	}

	return header, err
}

func (r *repository) GetZhkCharacteristics(ctx context.Context, zhk domain.Zhk) (domain.ZhkCharacteristics, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var characteristics domain.ZhkCharacteristics

	err := r.db.QueryRow(ctx, getZhkCharacteristicsSQL, zhk.ID).Scan(
		&characteristics.Class, &characteristics.Decoration,
		&characteristics.CeilingHeight.HighestHeight, &characteristics.CeilingHeight.LowestHeight,
		&characteristics.Floors.HighestFloor, &characteristics.Floors.LowestFloor,
		&characteristics.Square.HighestSquare, &characteristics.Square.LowestSquare,
	)

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getZhkCharacteristicsSQL,
		"params":    logger.LoggerFields{"id": zhk.ID},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetZhkInformation failed")
	} else {
		r.logger.WithFields(logFields).Info("GetZhkInformation succeeded")
	}

	return characteristics, err
}

func (r *repository) GetZhkApartments(ctx context.Context, zhk domain.Zhk) (domain.ZhkApartments, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var apartments domain.ZhkApartments

	rows, err := r.db.Query(ctx, getZhkApartmentsSQL, zhk.ID)
	if err != nil {
		return apartments, err
	}
	defer rows.Close()

	for rows.Next() {
		var appartment domain.ZhkApartment
		err = rows.Scan(
			&appartment.Rooms, &appartment.LowestPrice, &appartment.HighestPrice,
			&appartment.MinSquare, &appartment.Offers,
		)
		if err != nil {
			return apartments, err
		}
		apartments.Apartments = append(apartments.Apartments, appartment)
	}

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getZhkApartmentsSQL,
		"params":    logger.LoggerFields{"id": zhk.ID},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetZhkApartments failed")
	} else {
		r.logger.WithFields(logFields).Info("GetZhkApartments succeeded")
	}

	return apartments, err
}

func (r *repository) GetZhkReviews(ctx context.Context, zhk domain.Zhk) (domain.ZhkReviews, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var reviews domain.ZhkReviews

	rows, err := r.db.Query(ctx, getZhkReviewsSQL, zhk.ID)
	if err != nil {
		return reviews, err
	}
	defer rows.Close()

	for rows.Next() {
		var review domain.Review
		err = rows.Scan(
			&review.Avatar, &review.FirstName, &review.LastName,
			&review.Rating, &review.Text,
		)
		if err != nil {
			return reviews, err
		}
		reviews.Reviews = append(reviews.Reviews, review)
	}

	err = r.db.QueryRow(ctx, getZhkReviewsParamsSQL, zhk.ID).Scan(
		&reviews.Quantity, &reviews.TotalRating,
	)
	if err != nil {
		return reviews, err
	}

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getZhkReviewsSQL,
		"params":    logger.LoggerFields{"id": zhk.ID},
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetZhkReviews failed")
	} else {
		r.logger.WithFields(logFields).Info("GetZhkReviews succeeded")
	}

	return reviews, err
}

func (r *repository) GetAllZhk(ctx context.Context) ([]domain.Zhk, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var zhks []domain.Zhk

	rows, err := r.db.Query(ctx, getAllZhkSQL)
	if err != nil {
		return zhks, err
	}
	defer rows.Close()

	for rows.Next() {
		var zhk domain.Zhk
		err = rows.Scan(
			&zhk.ID, &zhk.ClassID, &zhk.Name, &zhk.Developer,
			&zhk.Phone, &zhk.Address, &zhk.Description,
		)
		if err != nil {
			return zhks, err
		}
		zhks = append(zhks, zhk)
	}

	logFields := logger.LoggerFields{
		"requestID": requestID,
		"query":     getAllZhkSQL,
		"success":   err == nil,
	}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetAllZhk failed")
	} else {
		r.logger.WithFields(logFields).Info("GetAllZhk succeeded")
	}

	return zhks, err

}

// endregion
