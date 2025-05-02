package domain

import (
	"time"
)

type OfferDetails struct {
	Pages  int         `json:"pages"`
	Offers []OfferInfo `json:"offers"`
}

type OfferInfo struct {
	Offer     Offer     `json:"offer"`
	OfferData OfferData `json:"offer_data"`
}

type OfferData struct {
	Images    []OfferImage `json:"offer_images"`
	Seller    OfferSeller  `json:"seller"`
	Metro     Metro        `json:"metro"`
	OfferStat OfferStat    `json:"offer_stat"`
}

type Offer struct {
	ID             int       `json:"id"`
	SellerID       int       `json:"seller_id"`
	OfferTypeID    int       `json:"offer_type_id"`
	MetroStationID *int      `json:"metro_station_id,omitempty"`
	RentTypeID     *int      `json:"rent_type_id,omitempty"`
	PurchaseTypeID *int      `json:"purchase_type_id,omitempty"`
	PropertyTypeID int       `json:"property_type_id"`
	StatusID       int       `json:"-"`
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
	Longitude      string    `json:"logitude"`
	Latitude       string    `json:"latitude"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

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
	Page           *int    `json:"page"`
}

type Metro struct {
	Id      int    `json:"station_id"`
	Color   string `json:"color"`
	Station string `json:"station"`
}

type OfferImage struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

type OfferSeller struct {
	FirstName string    `json:"seller_name"`
	LastName  string    `json:"seller_last_name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type OfferStat struct {
	LikesStat LikesStat `json:"likes_stat"`
	Views     *int      `json:"views"`
	// FavoutiteStat *FavoriteStat `json:"favourite_stat"`
}

type LikesStat struct {
	IsLiked bool `json:"is_liked"`
	Amount  int  `json:"amount"`
}

// type FavoriteStat struct {
// 	IsFavourited bool `json:"is_favorited"`
// 	Amount       *int  `json:"amount"`
// }

type LikeRequest struct {
	OfferId int `json:"offer_id"`
	UserId  int `json:"user_id"`
}

const OfferStatusDraft = 2
