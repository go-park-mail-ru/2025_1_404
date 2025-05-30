package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
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
	SELECT id, class_id, name, developer, phone_number, address, description, metro_station_id
	FROM kvartirum.HousingComplex
	WHERE id = $1;
	`
	getZhkHeaderSQL = `
	SELECT
		COALESCE(ARRAY_AGG(DISTINCT img.uuid) FILTER (WHERE img.uuid IS NOT NULL), '{}') AS images,
		COUNT (DISTINCT img.id) as images_size,
		COALESCE(MIN(o.price), 0) AS lowest_price,
		COALESCE(MAX(o.price), 0) AS highest_price
	FROM kvartirum.housingcomplex hc
	LEFT JOIN kvartirum.HousingComplexImages hci on hci.housing_complex_id = hc.id
	LEFT JOIN kvartirum.Image img on img.id = hci.image_id
	LEFT JOIN kvartirum.offer o ON o.complex_id = hc.id
	WHERE hc.id = $1
	GROUP BY hc.id;
	`

	getZhkCharacteristicsSQL = `
	SELECT
		hcc.name as class_name,
		COALESCE(MIN(o.ceiling_height), 0) AS lowest_ceiling_height,
		COALESCE(MAX(o.ceiling_height), 0) AS highest_ceiling_height,
		COALESCE(MIN(o.floor), 0) AS lowest_floor,
		COALESCE(MAX(o.floor), 0) AS highest_floor,
		COALESCE(MIN(o.area), 0) AS lowest_square,
		COALESCE(MAX(o.area), 0) AS highest_square
		FROM kvartirum.housingcomplex hc
		LEFT JOIN kvartirum.housingcomplexclass hcc ON hcc.id = hc.class_id
		LEFT JOIN kvartirum.offer o ON o.complex_id = hc.id
		WHERE hc.id = $1
		GROUP BY hcc.name;
	`

	getAllZhkSQL = `
	SELECT 
		id, class_id, name, developer, phone_number, address, description, metro_station_id
	FROM kvartirum.housingcomplex;
	`

	getZhkMetro = `
	SELECT 
		ms.name as station_name, ml.name as line_name, ml.color
		FROM kvartirum.MetroStation ms
		JOIN kvartirum.MetroLine ml ON ms.metro_line_id = ml.id
		WHERE ms.id = $1;
	`
)

func (r *zhkRepository) GetZhkByID(ctx context.Context, id int64) (domain.Zhk, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var zhk domain.Zhk
	err := r.db.QueryRow(ctx, getZhkByIDSQL, id).Scan(
		&zhk.ID, &zhk.ClassID, &zhk.Name, &zhk.Developer,
		&zhk.Phone, &zhk.Address, &zhk.Description, &zhk.MetroStationId,
	)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getZhkByIDSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("GetZhkByID")

	return zhk, err
}

func (r *zhkRepository) GetZhkHeader(ctx context.Context, zhk domain.Zhk) (domain.ZhkHeader, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	header := domain.ZhkHeader{Name: zhk.Name}

	err := r.db.QueryRow(ctx, getZhkHeaderSQL, zhk.ID).Scan(
		&header.Images, &header.ImagesSize, &header.LowestPrice, &header.HighestPrice,
	)

	logFields := logger.LoggerFields{"requestID": requestID, "query": getZhkHeaderSQL, "params": logger.LoggerFields{"id": zhk.ID}, "success": err == nil}

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
		&characteristics.Class,
		&characteristics.CeilingHeight.LowestHeight,
		&characteristics.CeilingHeight.HighestHeight,
		&characteristics.Floors.LowestFloor,
		&characteristics.Floors.HighestFloor,
		&characteristics.Square.LowestSquare,
		&characteristics.Square.HighestSquare,
	)

	logFields := logger.LoggerFields{"requestID": requestID, "query": getZhkCharacteristicsSQL, "params": logger.LoggerFields{"id": zhk.ID}, "success": err == nil}

	if err != nil {
		r.logger.WithFields(logFields).Error("GetZhkInformation failed")
	} else {
		r.logger.WithFields(logFields).Info("GetZhkInformation succeeded")
	}

	return characteristics, err
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
			&zhk.Phone, &zhk.Address, &zhk.Description, &zhk.MetroStationId,
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

func (r *zhkRepository) GetZhkMetro(ctx context.Context, id int64) (domain.ZhkMetro, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var zhkMetro domain.ZhkMetro

	err := r.db.QueryRow(ctx, getZhkMetro, id).Scan(
		&zhkMetro.Station, &zhkMetro.Line, &zhkMetro.Color,
	)

	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "stationID": id, "success": err == nil}).Info("SQL Query: GetZhkMetro")

	return zhkMetro, err
}
