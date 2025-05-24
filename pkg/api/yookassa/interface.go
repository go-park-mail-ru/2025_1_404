package yookassa

//go:generate mockgen -source interface.go -destination=mocks/interface.go -package=mocks

type YookassaRepo interface {
	CreatePayment(amount int, description string, redirectUri string) (*CreatePaymentResponse, error)
	GetPayment(paymentId string) (*CreatePaymentResponse, error)
}
