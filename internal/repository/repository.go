// internal/repository/repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/filestorage"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	// --- Offers ---
	CreateOffer(ctx context.Context, offer Offer) (int64, error)
	GetOfferByID(ctx context.Context, id int64) (Offer, error)
	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error)
	GetAllOffers(ctx context.Context) ([]Offer, error)
	UpdateOffer(ctx context.Context, offer Offer) error
	DeleteOffer(ctx context.Context, id int64) error
}

type repository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewRepository(db *pgxpool.Pool, logger logger.Logger) Repository {
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
	if u.Image == "" {
		imageID = nil
	} else {
		var id int64
		id, err := r.GetImageByUUID(ctx, u.Image)
		if err != nil {
			return updatedUser, fmt.Errorf("image not found")
		}
		imageID = id
	}

	requestID := ctx.Value(utils.RequestIDKey)

	// Обновляем юзера
	var id sql.NullInt64
	err := r.db.QueryRow(ctx, updateUserSQL,
		imageID, u.FirstName, u.LastName, u.Email, u.Password, u.ID,
	).Scan(&updatedUser.ID, &updatedUser.FirstName, &updatedUser.LastName,
		&updatedUser.Email, &id)

	// Получаем имя картинки по id картинки
	fileName, _:= r.GetImageByID(ctx, id)
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

func (r *repository) GetImageByUUID (ctx context.Context, uuid string) (int64, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var id int64
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

func (r *repository) GetImageByID (ctx context.Context, id sql.NullInt64) (string, error) {
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
	var id int64
	err := r.db.QueryRow(ctx, createOfferSQL,
		o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID,
		o.PurchaseTypeID, o.PropertyTypeID, o.StatusID, o.RenovationID,
		o.ComplexID, o.Price, o.Description, o.Floor, o.TotalFloors,
		o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
	).Scan(&id)
	return id, err
}

func (r *repository) GetOfferByID(ctx context.Context, id int64) (Offer, error) {
	var o Offer
	err := r.db.QueryRow(ctx, getOfferByIDSQL, id).Scan(
		&o.ID, &o.SellerID, &o.OfferTypeID, &o.MetroStationID, &o.RentTypeID,
		&o.PurchaseTypeID, &o.PropertyTypeID, &o.StatusID, &o.RenovationID,
		&o.ComplexID, &o.Price, &o.Description, &o.Floor, &o.TotalFloors,
		&o.Rooms, &o.Address, &o.Flat, &o.Area, &o.CeilingHeight,
		&o.CreatedAt, &o.UpdatedAt,
	)
	return o, err
}

func (r *repository) GetOffersBySellerID(ctx context.Context, sellerID int64) ([]Offer, error) {
	rows, err := r.db.Query(ctx, getOffersBySellerSQL, sellerID)
	if err != nil {
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
	return offers, nil
}

func (r *repository) GetAllOffers(ctx context.Context) ([]Offer, error) {
	rows, err := r.db.Query(ctx, getAllOffersSQL)
	if err != nil {
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
	return offers, nil
}

func (r *repository) UpdateOffer(ctx context.Context, o Offer) error {
	_, err := r.db.Exec(ctx, updateOfferSQL,
		o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
		o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID,
		o.Price, o.Description, o.Floor, o.TotalFloors, o.Rooms,
		o.Address, o.Flat, o.Area, o.CeilingHeight, o.ID,
	)
	return err
}

func (r *repository) DeleteOffer(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, deleteOfferSQL, id)
	return err
}

// endregion
