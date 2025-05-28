//go:generate easyjson -all

package domain

import (
	"time"
)

//easyjson:json
type OfferInfo struct {
	Offer     Offer     `json:"offer"`
	OfferData OfferData `json:"offer_data"`
}

//easyjson:json
type OfferPromotion struct {
	IsPromoted    bool       `json:"is_promoted"`
	PromotedUntil *time.Time `json:"promoted_until"`
}

//easyjson:json
type OfferData struct {
	Images         []OfferImage        `json:"offer_images"`
	Seller         OfferSeller         `json:"seller"`
	Metro          Metro               `json:"metro"`
	OfferStat      OfferStat           `json:"offer_stat"`
	Prices         []OfferPriceHistory `json:"offer_prices"`
	Promotion      *OfferPromotion     `json:"offer_promotion"`
	PromotionScore float32             `json:"-"`
}

//easyjson:json
type Offer struct {
	ID             int        `json:"id"`
	SellerID       int        `json:"seller_id"`
	OfferTypeID    int        `json:"offer_type_id"`
	MetroStationID *int       `json:"metro_station_id,omitempty"`
	RentTypeID     *int       `json:"rent_type_id,omitempty"`
	PurchaseTypeID *int       `json:"purchase_type_id,omitempty"`
	PropertyTypeID int        `json:"property_type_id"`
	StatusID       int        `json:"-"`
	RenovationID   int        `json:"renovation_id"`
	ComplexID      *int       `json:"complex_id,omitempty"`
	Price          int        `json:"price"`
	Description    *string    `json:"description,omitempty"`
	Floor          int        `json:"floor"`
	TotalFloors    int        `json:"total_floors"`
	Rooms          int        `json:"rooms"`
	Address        *string    `json:"address,omitempty"`
	Flat           int        `json:"flat"`
	Area           int        `json:"area"`
	CeilingHeight  int        `json:"ceiling_height"`
	Longitude      string     `json:"logitude"`
	Latitude       string     `json:"latitude"`
	Verified       bool       `json:"verified"`
	Comment        *string    `json:"comment,omitempty"`
	PromotesUntil  *time.Time `json:"promotes_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

//easyjson:json
type OfferFilter struct {
	MinArea        *int    `json:"min_area"`
	MaxArea        *int    `json:"max_area"`
	MinPrice       *int    `json:"min_price"`
	MaxPrice       *int    `json:"max_price"`
	Floor          *int    `json:"floor"`
	Rooms          *int    `json:"rooms"`
	Address        *string `json:"address"`
	RenovationID   *int    `json:"renovation_id"`
	PropertyTypeID *int    `json:"property_type_id"`
	PurchaseTypeID *int    `json:"purchase_type_id"`
	RentTypeID     *int    `json:"rent_type_id"`
	OfferTypeID    *int    `json:"offer_type_id"`
	NewBuilding    *bool   `json:"new_building"`
	SellerID       *int    `json:"seller_id"`
	OnlyMe         *bool   `json:"me"`
}

//easyjson:json
type Metro struct {
	Id      int    `json:"station_id"`
	Color   string `json:"color"`
	Station string `json:"station"`
}

//easyjson:json
type OfferImage struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

//easyjson:json
type OfferSeller struct {
	FirstName string    `json:"seller_name"`
	LastName  string    `json:"seller_last_name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

//easyjson:json
type OfferStat struct {
	LikesStat    LikesStat    `json:"likes_stat"`
	Views        *int         `json:"views"`
	FavoriteStat FavoriteStat `json:"favorite_stat"`
}

//easyjson:json
type LikesStat struct {
	IsLiked bool `json:"is_liked"`
	Amount  int  `json:"amount"`
}

//easyjson:json
type LikeRequest struct {
	OfferId int `json:"offer_id"`
	UserId  int `json:"user_id"`
}

const OfferStatusDraft = 2

//easyjson:json
type OfferPriceHistory struct {
	Price int       `json:"price"`
	Date  time.Time `json:"date"`
}

//easyjson:json
type FavoriteRequest struct {
	UserId  int `json:"user_id"`
	OfferId int `json:"offer_id"`
}

//easyjson:json
type FavoriteStat struct {
	IsFavorited bool `json:"is_favorited"`
	Amount      int  `json:"amount"`
}

//easyjson:json
type CreatePaymentRequest struct {
	Type int `json:"type"`
}

//easyjson:json
type CreatePaymentResponse struct {
	OfferId    int32  `json:"offer_id"`
	PaymentUri string `json:"payment_uri"`
}

//easyjson:json
type CheckPaymentResponse struct {
	OfferId  int  `json:"offer_id"`
	IsActive bool `json:"is_active"`
	IsPaid   bool `json:"is_paid"`
	Days     int  `json:"days"`
}

//easyjson:json
type PaymentPeriods struct {
	Days  int
	Price int
}

//easyjson:json
type OffersInfo []OfferInfo

//easyjson:json
type Stations []Metro

//easyjson:json
type OfferID struct {
	Id int `json:"id"`
}

//easyjson:json
type ImageID struct {
	ImageID int64 `json:"image_id"`
}
