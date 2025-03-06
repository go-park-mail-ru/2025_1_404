package offers

// Offer Структура объявления
type Offer struct {
	ID           int      `json:"id"`
	Seller       string   `json:"seller"`
	PropertyType string   `json:"property_type"` // квартира, дом, апартаменты
	OfferType    string   `json:"offer_type"`    // продажа, аренда
	PurchaseType string   `json:"purchase_type"` // вторичка, новостройка
	RentType     string   `json:"rent_type"`     // долгосрок, посуточно
	Address      string   `json:"address"`
	MetroLine    string   `json:"metro_line"`
	MetroStation string   `json:"metro_station"`
	Floor        int      `json:"floor"`
	TotalFloors  int      `json:"total_floors"`
	Area         float64  `json:"area"`
	Rooms        int      `json:"rooms"`
	Price        int      `json:"price"`
	Photos       []string `json:"photos"`
	Description  string   `json:"description"`
}
