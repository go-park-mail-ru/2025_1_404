package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/api/yookassa"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
)

type paymentUsecase struct {
	repo         payment.PaymentRepository
	yookassaRepo yookassa.YookassaRepo
	logger       logger.Logger
	cfg          *config.Config
}

func NewPaymentUsecase(repo payment.PaymentRepository, yookassaRepo yookassa.YookassaRepo, logger logger.Logger, cfg *config.Config) *paymentUsecase {
	return &paymentUsecase{repo: repo, yookassaRepo: yookassaRepo, logger: logger, cfg: cfg}
}

func (u *paymentUsecase) CheckType(paymentType int) (exits bool) {
	paymentPeriods := u.GetPaymentPeriods()
	_, ok := paymentPeriods[paymentType]
	return ok
}

func (u *paymentUsecase) ValidatePayment(ctx context.Context, paymentId int, offerId int) bool {
	payment, err := u.repo.GetPaymentById(ctx, paymentId)
	if err != nil {
		return false
	}
	return payment.OfferId == offerId
}

func (u *paymentUsecase) CheckPayment(ctx context.Context, paymentId int) (*domain.OfferPayment, error) {
	payment, err := u.repo.GetPaymentById(ctx, paymentId)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: get payment by id failed")
		return nil, err
	}
	if payment.IsActive {
		paymentResponse, err := u.yookassaRepo.GetPayment(payment.YookassaId)
		if err != nil {
			u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: get payment failed")
			return nil, err
		}
		if paymentResponse.Status != "pending" {
			payment.IsPaid = paymentResponse.Paid
			err = u.repo.UpdatePayment(ctx, payment)
			if err != nil {
				u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: update payment failed")
				return nil, err
			}
			err = u.repo.DeactivateAllPaymentsByOfferId(ctx, payment.OfferId)
			if err != nil {
				u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: deactivate payment by id failed")
				return nil, err
			}
		}
	}
	return payment, nil
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, request *domain.CreatePaymentRequest) (*domain.CreatePaymentResponse, error) {
	paymentPeriods := u.GetPaymentPeriods()
	paymentPeriod, ok := paymentPeriods[request.Type]
	if !ok {
		return nil, fmt.Errorf("неизвестный тариф")
	}
	offerPayment, err := u.repo.CreatePaymentForOfferId(ctx, int(request.OfferId), request.Type)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: create payment for offerid failed")
		return nil, err
	}
	paymentResponse, err := u.yookassaRepo.CreatePayment(
		paymentPeriod.Price,
		fmt.Sprintf("Объявление №%d. Продвижение %d дней", request.OfferId, paymentPeriod.Days),
		fmt.Sprintf("%s/payment/check/%d", u.cfg.App.BaseFrontendDir, offerPayment.Id),
	)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: create payment failed")
		return nil, err
	}
	offerPayment.YookassaId = paymentResponse.Id
	err = u.repo.UpdatePayment(ctx, offerPayment)
	if err != nil {
		u.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Warn("Payment usecase: update payment failed")
		return nil, err
	}

	return &domain.CreatePaymentResponse{
		OfferId:    request.OfferId,
		PaymentUri: paymentResponse.Confirmation.ConfirmationUri,
	}, nil
}

func (u *paymentUsecase) GetPaymentPeriods() map[int]domain.PaymentPeriods {
	return map[int]domain.PaymentPeriods{
		1: {
			Days:  3,
			Price: 500,
		},
		2: {
			Days:  7,
			Price: 3000,
		},
		3: {
			Days:  30,
			Price: 10000,
		},
	}
}
