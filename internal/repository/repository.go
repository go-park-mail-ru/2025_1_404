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
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateImage(ctx context.Context, file filestorage.FileUpload) error
	GetImageByID(ctx context.Context, id sql.NullInt64) (string, error)

	// --- Offers ---
	CreateOffer(ctx context.Context, offer Offer) (int64, error)
	GetOfferByID(ctx context.Context, id int64) (Offer, error)
	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error)
	GetAllOffers(ctx context.Context) ([]Offer, error)
	UpdateOffer(ctx context.Context, offer Offer) error
	DeleteOffer(ctx context.Context, id int64) error
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
		SELECT id, image_id, first_name, last_name, email, password,
			last_notification_id, token_version, created_at, updated_at
		FROM kvartirum.Users
		WHERE email = $1;
	`

	getUserByIDSQL = `
		SELECT id, image_id, first_name, last_name, email, password,
			last_notification_id, token_version, created_at, updated_at
		FROM kvartirum.Users
		WHERE id = $1;
	`

	updateUserSQL = `
		UPDATE kvartirum.Users
		SET image_id = $1, first_name = $2, last_name = $3, email = $4,
			password = $5
		WHERE id = $6
		RETURNING id, first_name, last_name, email, image_id;
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
)

func (r *repository) CreateUser(ctx context.Context, u User) (int64, error) {
	var id int64
	requestID := ctx.Value(utils.RequestIDKey)
	err := r.db.QueryRow(ctx, createUserSQL,
		u.ImageID, u.FirstName, u.LastName, u.Email, u.Password, u.TokenVersion,
	).Scan(&id)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     createUserSQL,
		"params": logger.LoggerFields{
			"name":          u.FirstName,
			"last_name":     u.LastName,
			"email":         u.Email,
			"token_version": u.TokenVersion,
			"image_id":      u.ImageID,
		},
		"success": err == nil,
	}).Info("SQL query CreateUser")

	return id, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User

	requestID := ctx.Value(utils.RequestIDKey)
	err := r.db.QueryRow(ctx, getUserByEmailSQL, email).Scan(
		&u.ID, &u.ImageID, &u.FirstName, &u.LastName, &u.Email, &u.Password,
		&u.LastNotificationID, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt,
	)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     getUserByEmailSQL,
		"params": logger.LoggerFields{
			"email": email,
		},
		"success": err == nil,
	}).Info("SQL query GetUserByEmail")

	return u, err
}

func (r *repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	var u User

	requestID := ctx.Value(utils.RequestIDKey)
	err := r.db.QueryRow(ctx, getUserByIDSQL, id).Scan(
		&u.ID, &u.ImageID, &u.FirstName, &u.LastName, &u.Email, &u.Password,
		&u.LastNotificationID, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt,
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
	var updatedUser domain.User

	// По имени картинки ищем id в БД
	var imageID interface{}
	imgID, err := r.GetImageByUUID(ctx, u.Image)
	if imgID.Valid {
		imageID = imgID.Int64
	} else {
		imageID = nil
	}

	requestID := ctx.Value(utils.RequestIDKey)

	// Обновляем юзера
	var id sql.NullInt64
	err = r.db.QueryRow(ctx, updateUserSQL,
		imageID, u.FirstName, u.LastName, u.Email, u.Password, u.ID,
	).Scan(&updatedUser.ID, &updatedUser.FirstName, &updatedUser.LastName,
		&updatedUser.Email, &id)

	// Получаем имя картинки по id картинки
	fileName, _ := r.GetImageByID(ctx, id)
	updatedUser.Image = fileName

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     updateUserSQL,
		"params": logger.LoggerFields{
			"name":       u.FirstName,
			"last_name":  u.LastName,
			"email":      u.Email,
			"image_path": u.Image,
		},
		"success": err == nil,
	}).Info("SQL query UpdateUser")

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
		"requestID": requestID,
		"query":     getImageByIDSQL,
		"params": logger.LoggerFields{
			"id": id,
		},
		"success": err == nil,
	})

	return fileName, err

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

// endregion
