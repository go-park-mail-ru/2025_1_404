package repository_test

import (
	"context"
	"testing"
	"time"
	_ "time"

	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/repository"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (repository.Repository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := repository.NewRepository(mock, logger.NewStub())
	return repo, mock
}

func TestRepository_CreateUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()

	user := repository.User{
		FirstName:    "Ivan",
		LastName:     "Petrov",
		Email:        "ivan@mail.ru",
		Password:     "hashed_pw",
		TokenVersion: 1,
		ImageID:      nil,
	}

	mock.ExpectQuery(`(?i)INSERT INTO kvartirum.Users`).
		WithArgs((*int64)(nil), user.FirstName, user.LastName, user.Email, user.Password, user.TokenVersion).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, int64(42), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByEmail(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	email := "user@example.com"
	id := int64(1)

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name, u.last_name, u.email, u.password\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.email = \$1`).
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", email, "hashed_pw"))

	u, err := repo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "avatar.png", u.Image)
	require.Equal(t, email, u.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByID(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)

	mock.ExpectQuery(`(?i)SELECT\s+u.id,\s+COALESCE\(i.uuid, ''\) as image,\s+u.first_name, u.last_name, u.email, u.password\s+FROM kvartirum.Users u\s+LEFT JOIN kvartirum.Image i on u.image_id = i.id\s+WHERE u.id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "image", "first_name", "last_name", "email", "password",
		}).AddRow(id, "avatar.png", "Ivan", "Petrov", "user@example.com", "hashed_pw"))

	u, err := repo.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, 1, u.ID)
	require.Equal(t, "Ivan", u.FirstName)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.Background()
	user := domain.User{
		ID:        1,
		FirstName: "Ivan",
		LastName:  "Petrov",
		Email:     "new@mail.ru",
		Image:     "new.png",
	}

	mock.ExpectQuery(`(?i)UPDATE kvartirum.Users`).
		WithArgs(user.Image, user.FirstName, user.LastName, user.Email, user.ID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "first_name", "last_name", "email", "image_uuid",
		}).AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Image))

	updated, err := repo.UpdateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, user.ID, updated.ID)
	require.Equal(t, user.Image, updated.Image)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteUser(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := int64(1)

	mock.ExpectExec(`(?i)DELETE FROM kvartirum.Users WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteUser(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateOffer(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	o := repository.Offer{
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

	o := repository.Offer{ID: 1, SellerID: 2, OfferTypeID: 1, PropertyTypeID: 1, StatusID: 1, RenovationID: 1, Price: 12345, Floor: 1, TotalFloors: 3, Rooms: 2, Flat: 10, Area: 40, CeilingHeight: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}

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

	o := repository.Offer{ID: 1, OfferTypeID: 1, PropertyTypeID: 1, StatusID: 1, RenovationID: 1, Price: 100000, Floor: 2, TotalFloors: 5, Rooms: 3, Flat: 12, Area: 40, CeilingHeight: 3}

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
