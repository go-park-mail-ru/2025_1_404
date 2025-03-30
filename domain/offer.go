package domain

import "time"

type Offer struct {
	ID             int       `json:"id"`
	SellerID       int       `json:"seller_id"`
	OfferTypeID    int       `json:"offer_type_id"`
	MetroStationID *int      `json:"metro_station_id,omitempty"`
	RentTypeID     *int      `json:"rent_type_id,omitempty"`
	PurchaseTypeID *int      `json:"purchase_type_id,omitempty"`
	PropertyTypeID int       `json:"property_type_id"`
	StatusID       int       `json:"status_id"`
	RenovationID   int       `json:"renovation_id"`
	ComplexID      *int      `json:"complex_id,omitempty"`
	Price          int       `json:"price"`
	Description    *string   `json:"description,omitempty"`
	Floor          int       `json:"floor"`
	TotalFloors    int       `json:"total_floors"`
	Rooms          int       `json:"rooms"`
	Address        *string   `json:"address,omitempty"`
	Flat           int       `json:"flat"`
	Area           int       `json:"area"`
	CeilingHeight  int       `json:"ceiling_height"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
