package repository

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	database "github.com/go-park-mail-ru/2025_1_404/pkg/database/postgres"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

type paymentRepository struct {
	db     database.DB
	logger logger.Logger
}

func NewPaymentRepository(db database.DB, logger logger.Logger) *paymentRepository {
	return &paymentRepository{db: db, logger: logger}
}

const (
	getPaymentByIDSQL = `
	SELECT id, offer_id, yookassa_id, type, is_active, is_paid
	FROM kvartirum.OfferPayment
	WHERE id = $1;
	`
	updatePaymentByIDSQL = `
	UPDATE kvartirum.OfferPayment SET
		is_active = $2, is_paid = $3, yookassa_id = $4
		WHERE id = $1
	`

	deactivateAllPaymentsByOfferIdSQL = `
	UPDATE kvartirum.OfferPayment SET
		is_active = false
		WHERE offer_id = $1 AND is_active = true;
	`

	createPaymentForOfferIdSQL = `
	INSERT INTO kvartirum.OfferPayment (offer_id, type, yookassa_id)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
)

func (r *paymentRepository) GetPaymentById(ctx context.Context, id int) (*domain.OfferPayment, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var payment domain.OfferPayment
	err := r.db.QueryRow(ctx, getPaymentByIDSQL, id).Scan(
		&payment.Id, &payment.OfferId, &payment.YookassaId, &payment.Type, &payment.IsActive, &payment.IsPaid,
	)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": getPaymentByIDSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("GetPaymentById")

	return &payment, err
}

func (r *paymentRepository) DeactivateAllPaymentsByOfferId(ctx context.Context, id int) error {
	requestID := ctx.Value(utils.RequestIDKey)

	_, err := r.db.Exec(ctx, deactivateAllPaymentsByOfferIdSQL, id)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": deactivateAllPaymentsByOfferIdSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("DeactivateAllPaymentsByOfferId")

	return err
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *domain.OfferPayment) error {
	requestID := ctx.Value(utils.RequestIDKey)
	_, err := r.db.Exec(ctx, updatePaymentByIDSQL, payment.Id, payment.IsActive, payment.IsPaid, payment.YookassaId)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": updatePaymentByIDSQL,
		"params": logger.LoggerFields{"id": payment.Id}, "success": err == nil}).Info("UpdatePayment")

	return err
}

func (r *paymentRepository) CreatePaymentForOfferId(ctx context.Context, id int, paymentType int) (*domain.OfferPayment, error) {
	requestID := ctx.Value(utils.RequestIDKey)

	var payment domain.OfferPayment
	payment.OfferId = id
	payment.Type = paymentType
	payment.IsActive = true
	payment.IsPaid = false
	payment.YookassaId = ""
	err := r.db.QueryRow(ctx, createPaymentForOfferIdSQL, id, paymentType, "").Scan(
		&payment.Id,
	)
	r.logger.WithFields(logger.LoggerFields{"requestID": requestID, "query": createPaymentForOfferIdSQL,
		"params": logger.LoggerFields{"id": id}, "success": err == nil}).Info("CreatePaymentForOfferId")

	return &payment, err
}
