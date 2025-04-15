package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (ZhkRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewZhkRepository(mock, logger.NewStub())
	return repo, mock
}

func TestRepository_GetZhkByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)
	expected := domain.Zhk{
		ID:          id,
		ClassID:     2,
		Name:        "ЖК Лесной",
		Developer:   "Брусника",
		Phone:       "8001234567",
		Address:     "Москва, Лесная 7",
		Description: "Уютный ЖК у парка",
	}

	mock.ExpectQuery(`(?i)SELECT id, class_id, name, developer, phone_number, address, description FROM kvartirum.HousingComplex WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "class_id", "name", "developer", "phone_number", "address", "description",
		}).AddRow(expected.ID, expected.ClassID, expected.Name, expected.Developer, expected.Phone, expected.Address, expected.Description))

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

	mock.ExpectQuery(`(?i)SELECT.*MIN\(o.price\).*MAX\(o.price\).*kvartirum.housingcomplex hc`).
		WithArgs(zhk.ID).
		WillReturnRows(pgxmock.NewRows([]string{"lowest_price", "highest_price", "images", "images_size"}).
			AddRow(expected.LowestPrice, expected.HighestPrice, expected.Images, expected.ImagesSize))

	got, err := repo.GetZhkHeader(context.Background(), zhk)
	require.NoError(t, err)
	require.Equal(t, expected.LowestPrice, got.LowestPrice)
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
			"class_name", "decoration", "max_ceiling_height", "min_ceiling_height",
			"max_floors", "min_floors", "max_area", "min_area",
		}).AddRow("Комфорт", []string{"Чистовая", "Без отделки"}, 3, 2, 25, 10, 120, 40))

	got, err := repo.GetZhkCharacteristics(context.Background(), zhk)
	require.NoError(t, err)
	require.Equal(t, "Комфорт", got.Class)
	require.Equal(t, 25, got.Floors.HighestFloor)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetZhkApartments(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	zhk := domain.Zhk{ID: 1}

	mock.ExpectQuery(`(?i)SELECT.*FROM kvartirum.offer o.*WHERE o.complex_id = \$1`).
		WithArgs(zhk.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"rooms", "lowest_price", "highest_price", "min_square", "offers",
		}).AddRow(1, 3000000, 4000000, 35, 10))

	got, err := repo.GetZhkApartments(context.Background(), zhk)
	require.NoError(t, err)
	require.Len(t, got.Apartments, 1)
	require.Equal(t, 1, got.Apartments[0].Rooms)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllZhk(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT.*FROM kvartirum.housingcomplex`).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "class_id", "name", "developer", "phone_number", "address", "description",
		}).AddRow(1, 2, "ЖК Радуга", "ПИК", "88005553535", "г. Москва", "Описание ЖК"))

	got, err := repo.GetAllZhk(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "ЖК Радуга", got[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}
