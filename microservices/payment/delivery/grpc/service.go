package service

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/microservices/payment"
	domain "github.com/go-park-mail-ru/2025_1_404/microservices/payment/domain"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	paymentpb "github.com/go-park-mail-ru/2025_1_404/proto/payment"
)

type paymentService struct {
	UC     payment.PaymentUsecase
	logger logger.Logger
	paymentpb.UnimplementedPaymentServiceServer
}

func NewPaymentService(usecase payment.PaymentUsecase, logger logger.Logger) *paymentService {
	return &paymentService{UC: usecase, logger: logger, UnimplementedPaymentServiceServer: paymentpb.UnimplementedPaymentServiceServer{}}
}

func (s *paymentService) CheckType(ctx context.Context, request *paymentpb.CheckTypeRequest) (*paymentpb.CheckTypeResponse, error) {
	exists := s.UC.CheckType(int(request.Type))
	return &paymentpb.CheckTypeResponse{
		IsValid: exists,
	}, nil
}

func (s *paymentService) CreatePayment(ctx context.Context, request *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	paymentResponse, err := s.UC.CreatePayment(ctx, &domain.CreatePaymentRequest{
		OfferId: request.OfferId,
		Type:    int(request.Type),
	})
	if err != nil {
		s.logger.Warn("failed to create payment")
		return nil, err
	}
	return &paymentpb.CreatePaymentResponse{
		OfferId:     paymentResponse.OfferId,
		RedirectUri: paymentResponse.PaymentUri,
	}, nil
}

func (s *paymentService) CheckPayment(ctx context.Context, request *paymentpb.CheckPaymentRequest) (*paymentpb.CheckPaymentResponse, error) {
	payment, err := s.UC.CheckPayment(ctx, int(request.PaymentId))
	if err != nil {
		s.logger.Warn("failed to check payment")
		return nil, err
	}
	return &paymentpb.CheckPaymentResponse{
		OfferId:  int32(payment.OfferId),
		IsActive: payment.IsActive,
		IsPaid:   payment.IsPaid,
		Days:     int32(payment.Days),
	}, nil
}

func (s *paymentService) ValidatePayment(ctx context.Context, request *paymentpb.ValidatePaymentRequest) (*paymentpb.ValidatePaymentResponse, error) {
	validationResult := s.UC.ValidatePayment(ctx, int(request.PaymentId), int(request.OfferId))
	return &paymentpb.ValidatePaymentResponse{
		IsValid: validationResult,
	}, nil
}
