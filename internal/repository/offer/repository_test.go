package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*offerRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewOfferRepository(mock, logger.NewStub())
	return repo, mock
}

func TestRepository_CreateOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	o := Offer{
		SellerID:       1,
		OfferTypeID:    1,
		PropertyTypeID: 1,
		StatusID:       1,
		RenovationID:   1,
		Price:          1000000,
		Floor:          2,
		TotalFloors:    5,
		Rooms:          2,
		Flat:           10,
		Area:           45,
		CeilingHeight:  3,
	}

	mock.ExpectQuery(`(?i)INSERT INTO kvartirum.Offer`).
		WithArgs(
			o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID,
			o.PurchaseTypeID, o.PropertyTypeID, o.StatusID, o.RenovationID,
			o.ComplexID, o.Price, o.Description, o.Floor, o.TotalFloors,
			o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
		).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))

	id, err := repo.CreateOffer(context.Background(), o)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetOfferByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	o := Offer{ID: 1, SellerID: 2, OfferTypeID: 1, PropertyTypeID: 1, StatusID: 1, RenovationID: 1, Price: 12345, Floor: 1, TotalFloors: 3, Rooms: 2, Flat: 10, Area: 40, CeilingHeight: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mock.ExpectQuery(`(?i)SELECT .* FROM kvartirum.Offer WHERE id = \$1`).
		WithArgs(o.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "seller_id", "offer_type_id", "metro_station_id", "rent_type_id", "purchase_type_id",
			"property_type_id", "offer_status_id", "renovation_id", "complex_id", "price", "description",
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height", "created_at", "updated_at",
		}).AddRow(
			o.ID, o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
			o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID, o.Price, o.Description,
			o.Floor, o.TotalFloors, o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
			o.CreatedAt, o.UpdatedAt,
		))

	got, err := repo.GetOfferByID(context.Background(), o.ID)
	require.NoError(t, err)
	require.Equal(t, o.ID, got.ID)
	require.Equal(t, o.Price, got.Price)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllOffers(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT .* FROM kvartirum.Offer`).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "seller_id", "offer_type_id", "metro_station_id", "rent_type_id", "purchase_type_id",
			"property_type_id", "offer_status_id", "renovation_id", "complex_id", "price", "description",
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height", "created_at", "updated_at",
		}).AddRow(1, 2, 1, nil, nil, nil, 1, 1, 1, nil, 100000, nil, 2, 5, 2, nil, 10, 50, 3, time.Now(), time.Now()))

	list, err := repo.GetAllOffers(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	o := Offer{ID: 1, OfferTypeID: 1, PropertyTypeID: 1, StatusID: 1, RenovationID: 1, Price: 100000, Floor: 2, TotalFloors: 5, Rooms: 3, Flat: 12, Area: 40, CeilingHeight: 3}

	mock.ExpectExec(`(?i)UPDATE kvartirum.Offer`).
		WithArgs(
			o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
			o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID,
			o.Price, o.Description, o.Floor, o.TotalFloors, o.Rooms,
			o.Address, o.Flat, o.Area, o.CeilingHeight, o.ID,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.UpdateOffer(context.Background(), o)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.Offer WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteOffer(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetOffersByFilter(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	filter := domain.OfferFilter{
		MinArea:  ptr(40),
		MaxPrice: ptr(2000000),
	}

	mock.ExpectQuery(`(?i)SELECT .* FROM kvartirum.Offer WHERE area >= \$1 AND price <= \$2;`).
		WithArgs(*filter.MinArea, *filter.MaxPrice).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "seller_id", "offer_type_id", "metro_station_id", "rent_type_id", "purchase_type_id",
			"property_type_id", "offer_status_id", "renovation_id", "complex_id", "price", "description",
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height", "created_at", "updated_at",
		}).AddRow(1, 2, 1, nil, nil, nil, 1, 1, 1, nil, 1800000, nil, 2, 5, 2, nil, 10, 50, 3, time.Now(), time.Now()))

	offers, err := repo.GetOffersByFilter(context.Background(), filter)
	require.NoError(t, err)
	require.NotEmpty(t, offers)
	require.NoError(t, mock.ExpectationsWereMet())
}

func ptr[T any](v T) *T {
	return &v
}

func TestRepository_CreateImageAndBindToOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	offerID := 1
	uuid := "image-uuid"

	// Ожидаем вставку картинки
	mock.ExpectQuery(`(?i)INSERT INTO kvartirum.Image \(uuid\)`).
		WithArgs(uuid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(10)))

	// Ожидаем связь с оффером
	mock.ExpectExec(`(?i)INSERT INTO kvartirum.OfferImages`).
		WithArgs(offerID, int64(10)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	imageID, err := repo.CreateImageAndBindToOffer(ctx, offerID, uuid)
	require.NoError(t, err)
	require.Equal(t, int64(10), imageID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteOfferImage(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	imageID := int64(5)

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.OfferImages`).
		WithArgs(imageID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteOfferImage(context.Background(), imageID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetOfferImageWithUUID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	imageID := int64(7)
	expectedOfferID := int64(3)
	expectedUUID := "abc-uuid"

	mock.ExpectQuery(`(?i)SELECT oi.offer_id, i.uuid FROM kvartirum.OfferImages oi`).
		WithArgs(imageID).
		WillReturnRows(pgxmock.NewRows([]string{"offer_id", "uuid"}).AddRow(expectedOfferID, expectedUUID))

	offerID, uuid, err := repo.GetOfferImageWithUUID(context.Background(), imageID)
	require.NoError(t, err)
	require.Equal(t, expectedOfferID, offerID)
	require.Equal(t, expectedUUID, uuid)
	require.NoError(t, mock.ExpectationsWereMet())
}
