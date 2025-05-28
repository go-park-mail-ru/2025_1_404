package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/microservices/zhk/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/lib/pq"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*zhkRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewZhkRepository(mock, logger.NewStub())
	return repo, mock
}

func TestRepository_GetZhkByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)
	station := 3
	expected := domain.Zhk{
		ID:             id,
		ClassID:        2,
		Name:           "ЖК Лесной",
		Developer:      "Брусника",
		Phone:          "8001234567",
		Address:        "Москва, Лесная 7",
		Description:    "Уютный ЖК у парка",
		MetroStationId: &station,
	}

	mock.ExpectQuery(`(?i)SELECT id, class_id, name, developer, phone_number, address, description, metro_station_id FROM kvartirum.HousingComplex WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "class_id", "name", "developer", "phone_number", "address", "description", "metro_station_id",
		}).AddRow(expected.ID, expected.ClassID, expected.Name, expected.Developer, expected.Phone, expected.Address, expected.Description, expected.MetroStationId))
	got, err := repo.GetZhkByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetZhkHeader(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	zhk := domain.Zhk{ID: 1, Name: "ЖК Альфа"}
	expected := domain.ZhkHeader{
		Name:         zhk.Name,
		LowestPrice:  3000000,
		HighestPrice: 6000000,
		Images:       []string{"img1", "img2"},
		ImagesSize:   2,
	}

	mock.ExpectQuery(`(?i)SELECT COALESCE\(ARRAY_AGG\(DISTINCT img.uuid\).*images_size`).
		WithArgs(zhk.ID).
		WillReturnRows(pgxmock.NewRows([]string{"images", "images_size"}).
			AddRow(pq.StringArray{"img1", "img2"}, expected.ImagesSize))

	got, err := repo.GetZhkHeader(context.Background(), zhk)
	require.NoError(t, err)
	require.Equal(t, expected.Images, got.Images)
	require.Equal(t, expected.ImagesSize, got.ImagesSize)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetZhkCharacteristics(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	zhk := domain.Zhk{ID: 1}

	mock.ExpectQuery(`(?i)SELECT.*FROM kvartirum.housingcomplex hc`).
		WithArgs(zhk.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"class_name",
		}).AddRow("Комфорт"))

	got, err := repo.GetZhkCharacteristics(context.Background(), zhk)
	require.NoError(t, err)
	require.Equal(t, "Комфорт", got.Class)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllZhk(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT.*FROM kvartirum.housingcomplex`).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "class_id", "name", "developer", "phone_number", "address", "description", "metro_station_id",
		}).AddRow(1, 2, "ЖК Радуга", "ПИК", "88005553535", "г. Москва", "Описание ЖК", nil))

	got, err := repo.GetAllZhk(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "ЖК Радуга", got[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetZhkMetro(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	stationID := int64(3)

	expected := domain.ZhkMetro{
		Station: "Кропоткинская",
		Line:    "Красная",
		Color:   "#FF0000",
	}

	mock.ExpectQuery(`(?i)SELECT\s+ms.name as station_name, ml.name as line_name, ml.color\s+FROM kvartirum.MetroStation ms\s+JOIN kvartirum.MetroLine ml ON ms.metro_line_id = ml.id\s+WHERE ms.id = \$1`).
		WithArgs(stationID).
		WillReturnRows(pgxmock.NewRows([]string{"station_name", "line_name", "color"}).
			AddRow(expected.Station, expected.Line, expected.Color))

	result, err := repo.GetZhkMetro(ctx, stationID)
	require.NoError(t, err)
	require.Equal(t, expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
