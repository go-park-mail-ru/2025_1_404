package payment

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_zhk_repo.go -package=mocks

type PaymentRepository interface {
	GetPaymentById(ctx context.Context, id int) (*domain.OfferPayment, error)
	UpdatePayment(ctx context.Context, payment *domain.OfferPayment) error
	DeactivateAllPaymentsByOfferId(ctx context.Context, id int) error
	CreatePaymentForOfferId(ctx context.Context, id int, paymentType int) (*domain.OfferPayment, error)
}
