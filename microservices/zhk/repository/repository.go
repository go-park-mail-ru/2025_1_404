package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type zhkRepository struct {
	db     database.DB
	logger logger.Logger
}

func NewZhkRepository(db database.DB, logger logger.Logger) *zhkRepository {
	return &zhkRepository{db: db, logger: logger}
}

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

func (r *zhkRepository) GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var zhk domain.Zhk
	err := r.db.QueryRow(ctx, getZhkByIDSQL, id).Scan(
		&zhk.ID, &zhk.ClassID, &zhk.Name, &zhk.Developer,
		&zhk.Phone, &zhk.Address, &zhk.Description,
	)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getZhkByIDSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("GetZhkByID")

	return zhk, err
}

func (r *zhkRepository) GetZhkHeader(ctx context.Context, zhk domain.Zhk) (domain.ZhkHeader, error) {
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

func (r *zhkRepository) GetZhkCharacteristics(ctx context.Context, zhk domain.Zhk) (domain.ZhkCharacteristics, error) {
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

func (r *zhkRepository) GetZhkApartments(ctx context.Context, zhk domain.Zhk) (domain.ZhkApartments, error) {
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

func (r *zhkRepository) GetZhkReviews(ctx context.Context, zhk domain.Zhk) (domain.ZhkReviews, error) {
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

func (r *zhkRepository) GetAllZhk(ctx context.Context) ([]domain.Zhk, error) {
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
