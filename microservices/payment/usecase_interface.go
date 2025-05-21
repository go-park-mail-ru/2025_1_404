package payment

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_ai.go -package=mocks

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, request *domain.CreatePaymentRequest) (*domain.CreatePaymentResponse, error)
	CheckType(paymentType int) bool
	CheckPayment(ctx context.Context, paymentId int) (*domain.OfferPayment, error)
	ValidatePayment(ctx context.Context, paymentId int, offerId int) bool
}
