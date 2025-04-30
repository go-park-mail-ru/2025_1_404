package domain

type Offer struct {
	OfferType     string  `json:"offer_type"`
	MetroStation  *string `json:"metro_station"`
	RentType      *string `json:"rent_type"`
	PurchaseType  string  `json:"purchase_type"`
	PropertyType  string  `json:"property_type"`
	Renovation    string  `json:"renovation"`
	Complex       *string `json:"complex_id"`
	Floor         int     `json:"floor"`
	TotalFloors   int     `json:"total_floors"`
	Rooms         int     `json:"rooms"`
	Address       string  `json:"address"`
	Area          int     `json:"area"`
	CeilingHeight int     `json:"ceiling_height"`
}

type MarketPrice struct {
	Total          float32 `json:"total"`
	PerSquareMeter float32 `json:"per_square_meter"`
}

type PossibleCostRange struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type EvaluationResult struct {
	MarketPrice       MarketPrice       `json:"market_price"`
	PossibleCostRange PossibleCostRange `json:"possible_cost_range"`
}
