package domain

type CreatePaymentRequest struct {
	OfferId int32 `json:"offer_id"`
	Type    int   `json:"type"`
}

type CreatePaymentResponse struct {
	OfferId    int32  `json:"offer_id"`
	PaymentUri string `json:"payment_uri"`
}

type PaymentPeriods struct {
	Days  int
	Price int
}

type OfferPayment struct {
	Id         int
	OfferId    int
	YookassaId string
	Type       int
	IsActive   bool
	IsPaid     bool
	Days       int
}
