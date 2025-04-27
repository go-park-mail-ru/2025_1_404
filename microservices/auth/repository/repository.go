package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/auth/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
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

type authRepository struct {
	db     database.DB
	logger logger.Logger
}

func NewAuthRepository(db database.DB, logger logger.Logger) *authRepository {
	return &authRepository{db: db, logger: logger}
}

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

func (r *authRepository) CreateUser(ctx context.Context, u User) (int64, error) {
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

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var u domain.User
	err := r.db.QueryRow(ctx, getUserByEmailSQL, email).Scan(
		&u.ID, &u.Image, &u.FirstName, &u.LastName, &u.Email, &u.Password,
	)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getUserByEmailSQL,
		"params": logger.LoggerFields{"email": email}, "success": err == nil}).Info("SQL query GetUserByEmail")

	return u, err
}

func (r *authRepository) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
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

func (r *authRepository) UpdateUser(ctx context.Context, u domain.User) (domain.User, error) {
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

func (r *authRepository) DeleteUser(ctx context.Context, id int64) error {
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

func (r *authRepository) CreateImage(ctx context.Context, fileName string) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, createImageSQL, fileName)

	r.logger.WithFields(logger.LoggerFields{
		"requestID": requestID,
		"query":     createImageSQL,
		"params": logger.LoggerFields{
			"name": fileName,
		},
		"success": err == nil,
	})

	return err
}

func (r *authRepository) GetImageByUUID(ctx context.Context, uuid string) (sql.NullInt64, error) {
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

func (r *authRepository) GetImageByID(ctx context.Context, id sql.NullInt64) (string, error) {
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

func (r *authRepository) DeleteUserImage(ctx context.Context, id int64) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, deleteUserImageSQL, id)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": deleteUserSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("SQL DeleteImage")

	return err
}
