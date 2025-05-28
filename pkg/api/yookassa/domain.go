package yookassa

type CreatePaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Capture      bool         `json:"capture"`
	Confirmation Confirmation `json:"confirmation"`
	Description  string       `json:"description"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type            string `json:"type"`
	ReturnUri       string `json:"return_url,omitempty"`
	ConfirmationUri string `json:"confirmation_url,omitempty"`
}

type CreatePaymentResponse struct {
	Id           string       `json:"id"`
	Status       string       `json:"status"`
	Paid         bool         `json:"paid"`
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	CreatedAt    string       `json:"created_at"`
	Description  string       `json:"description"`
	Metadata     any          `json:"metadata"`
	Recipient    Recipient    `json:"recipient"`
	Refundable   bool         `json:"refundable"`
	Test         bool         `json:"test"`
}

type Recipient struct {
	AccountId string `json:"account_id"`
	GatewayId string `json:"gateway_id"`
}
