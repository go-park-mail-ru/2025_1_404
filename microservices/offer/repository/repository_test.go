package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/microservices/offer/domain"
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
		Longitude:      "37.6173",
		Latitude:       "55.7558",
	}

	mock.ExpectQuery(`(?i)INSERT INTO kvartirum.Offer`).
		WithArgs(
			o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID,
			o.PurchaseTypeID, o.PropertyTypeID, o.StatusID, o.RenovationID,
			o.ComplexID, o.Price, o.Description, o.Floor, o.TotalFloors,
			o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
			o.Longitude, o.Latitude,
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

	o := Offer{
		ID:             1,
		SellerID:       2,
		OfferTypeID:    1,
		PropertyTypeID: 1,
		StatusID:       1,
		RenovationID:   1,
		Price:          12345,
		Floor:          1,
		TotalFloors:    3,
		Rooms:          2,
		Flat:           10,
		Area:           40,
		CeilingHeight:  3,
		Longitude:      "37.6173",
		Latitude:       "55.7558",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectQuery(`(?i)SELECT .* FROM kvartirum.Offer WHERE id = \$1`).
		WithArgs(o.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "seller_id", "offer_type_id", "metro_station_id", "rent_type_id", "purchase_type_id",
			"property_type_id", "offer_status_id", "renovation_id", "complex_id", "price", "description",
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height",
			"longitude", "latitude", "created_at", "updated_at",
		}).AddRow(
			o.ID, o.SellerID, o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
			o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID, o.Price, o.Description,
			o.Floor, o.TotalFloors, o.Rooms, o.Address, o.Flat, o.Area, o.CeilingHeight,
			o.Longitude, o.Latitude, o.CreatedAt, o.UpdatedAt,
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
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height",
			"longitude", "latitude", "created_at", "updated_at",
		}).AddRow(
			1, 2, 1, nil, nil, nil, 1, 1, 1, nil, 100000, nil,
			2, 5, 2, nil, 10, 50, 3,
			"37.6173", "55.7558", time.Now(), time.Now(),
		))

	list, err := repo.GetAllOffers(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	o := Offer{
		ID:             1,
		OfferTypeID:    1,
		MetroStationID: ptr(1),
		RentTypeID:     ptr(1),
		PurchaseTypeID: ptr(1),
		PropertyTypeID: 1,
		StatusID:       1,
		RenovationID:   1,
		ComplexID:      ptr(1),
		Price:          100000,
		Description:    ptr("Описание"),
		Floor:          2,
		TotalFloors:    5,
		Rooms:          3,
		Address:        ptr("Адрес"),
		Flat:           12,
		Area:           40,
		CeilingHeight:  3,
	}

	mock.ExpectExec(`(?i)UPDATE kvartirum.Offer`).
		WithArgs(
			o.OfferTypeID, o.MetroStationID, o.RentTypeID, o.PurchaseTypeID,
			o.PropertyTypeID, o.StatusID, o.RenovationID, o.ComplexID,
			o.Price, o.Description, o.Floor, o.TotalFloors, o.Rooms,
			o.Address, o.Flat, o.Area, o.CeilingHeight, o.Longitude, o.Latitude, o.ID,
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

	timeNow := time.Now()
	mock.ExpectQuery(`(?i)SELECT id, seller_id.*FROM kvartirum.Offer WHERE area >= \$1 AND price <= \$2 AND offer_status_id = \$3;`).
		WithArgs(*filter.MinArea, *filter.MaxPrice, 1).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "seller_id", "offer_type_id", "metro_station_id", "rent_type_id", "purchase_type_id",
			"property_type_id", "offer_status_id", "renovation_id", "complex_id", "price", "description",
			"floor", "total_floors", "rooms", "address", "flat", "area", "ceiling_height", "longitude", "latitude", "created_at", "updated_at",
		}).AddRow(
			1, 2, 1, nil, nil, nil, 1, 1, 1, nil, 1800000, nil, 2, 5, 2, nil, 10, 50, 3, "37.6173", "55.7558", timeNow, timeNow,
		))

	offers, err := repo.GetOffersByFilter(context.Background(), filter, nil)
	require.NoError(t, err)
	require.Len(t, offers, 1)
	require.Equal(t, int64(1), offers[0].ID)
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

func TestRepository_UpdateOfferStatus(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)UPDATE kvartirum.Offer SET offer_status_id = \$1 WHERE id = \$2;`).
		WithArgs(2, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.UpdateOfferStatus(context.Background(), 1, 2)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetStations(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT ms.id as station_id`).
		WillReturnRows(pgxmock.NewRows([]string{"station_id", "station_name", "color"}).AddRow(1, "Test Station", "Red"))

	stations, err := repo.GetStations(context.Background())
	require.NoError(t, err)
	require.Len(t, stations, 1)
	require.Equal(t, "Test Station", stations[0].Station)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IsOfferLiked(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT EXISTS`).
		WithArgs(42, 77).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	liked, err := repo.IsOfferLiked(context.Background(), domain.LikeRequest{UserId: 42, OfferId: 77})
	require.NoError(t, err)
	require.True(t, liked)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateLike(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)INSERT INTO kvartirum.Likes`).
		WithArgs(42, 77).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateLike(context.Background(), domain.LikeRequest{UserId: 42, OfferId: 77})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteLike(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.Likes`).
		WithArgs(42, 77).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteLike(context.Background(), domain.LikeRequest{UserId: 42, OfferId: 77})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetLikeStat(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`(?i)SELECT COUNT\(\*\) FROM kvartirum.Likes`).
		WithArgs(ptr(77)).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.GetLikeStat(context.Background(), domain.LikeRequest{OfferId: 77})
	require.NoError(t, err)
	require.Equal(t, 5, count)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IncrementView(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)INSERT INTO kvartirum.Views`).
		WithArgs(77).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.IncrementView(context.Background(), 77)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddOrUpdatePriceHistory(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)INSERT INTO kvartirum.OfferPriceHistory`).
		WithArgs(int64(1), 123456).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.AddOrUpdatePriceHistory(context.Background(), 1, 123456)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePriceHistory(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.OfferPriceHistory`).
		WithArgs(int64(1)).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeletePriceHistory(context.Background(), 1)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetPriceHistory(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	timestamp := time.Now()
	mock.ExpectQuery(`(?i)SELECT price, recorded_at FROM kvartirum.OfferPriceHistory`).
		WithArgs(int64(1), 5).
		WillReturnRows(pgxmock.NewRows([]string{"price", "recorded_at"}).AddRow(123456, timestamp))

	history, err := repo.GetPriceHistory(context.Background(), 1, 5)
	require.NoError(t, err)
	require.Len(t, history, 1)
	require.Equal(t, 123456, history[0].Price)
	require.Equal(t, timestamp, history[0].Date)
	require.NoError(t, mock.ExpectationsWereMet())
}
