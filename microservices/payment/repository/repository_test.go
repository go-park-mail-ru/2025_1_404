package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*paymentRepository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewPaymentRepository(mock, logger.NewStub())
	return repo, mock
}

func TestGetPaymentById(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "request-id")

	payment := domain.OfferPayment {
		Id: 1, 
		OfferId: 2,
		YookassaId: "id",
		Type: 1,
		IsActive: true,
		IsPaid: true,
	}

	mock.ExpectQuery(`(?i)SELECT id, offer_id, yookassa_id, type, is_active, is_paid FROM kvartirum\.OfferPayment WHERE id = \$1;`).
		WithArgs(payment.Id).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "offer_id", "yookassa_id", "type", "is_active", "is_paid", 
		}).AddRow(payment.Id, payment.OfferId, payment.YookassaId, payment.Type, payment.IsActive, payment.IsPaid))

	resp, err := repo.GetPaymentById(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, payment, *resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeactivateAllPaymentsByOfferId(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "request-id")

	offerID := 1

	mock.ExpectExec(`(?i)UPDATE\s+kvartirum\.OfferPayment\s+SET\s+is_active\s+=\s+false\s+WHERE\s+offer_id\s+=\s+\$1\s+AND\s+is_active\s+=\s+true;`).
		WithArgs(offerID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1)) // предполагаем, что 1 строка была обновлена

	err := repo.DeactivateAllPaymentsByOfferId(ctx, offerID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "request-id")

	payment := &domain.OfferPayment{
		Id:          1,
		IsActive:    false,
		IsPaid:      true,
		YookassaId:  "yookassa-id-123",
	}

	mock.ExpectExec(`(?i)UPDATE\s+kvartirum\.OfferPayment\s+SET\s+is_active\s+=\s+\$2,\s+is_paid\s+=\s+\$3,\s+yookassa_id\s+=\s+\$4\s+WHERE\s+id\s+=\s+\$1`).
		WithArgs(payment.Id, payment.IsActive, payment.IsPaid, payment.YookassaId).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.UpdatePayment(ctx, payment)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}


func TestCreatePaymentForOfferId(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	ctx := context.WithValue(context.Background(), utils.RequestIDKey, "request-id")

	offerID := 5
	paymentType := 2
	newPaymentID := 10

	// Ожидаем SQL-запрос с учетом структуры INSERT и RETURNING
	mock.ExpectQuery(`(?i)INSERT INTO kvartirum\.OfferPayment \(offer_id, type, yookassa_id\) VALUES \(\$1, \$2, \$3\) RETURNING id;`).
		WithArgs(offerID, paymentType, "").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newPaymentID))

	result, err := repo.CreatePaymentForOfferId(ctx, offerID, paymentType)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, newPaymentID, result.Id)
	require.Equal(t, offerID, result.OfferId)
	require.Equal(t, paymentType, result.Type)
	require.True(t, result.IsActive)
	require.False(t, result.IsPaid)
	require.Equal(t, "", result.YookassaId)

	require.NoError(t, mock.ExpectationsWereMet())
}
