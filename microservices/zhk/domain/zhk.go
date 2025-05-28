//go:generate easyjson -all

package domain

// Zhk Структура ЖК
type Zhk struct {
	ID             int64
	ClassID        int64
	Name           string
	Developer      string
	Phone          string
	Address        string
	Description    string
	MetroStationId *int
}

//easyjson:json
type ZhkInfo struct {
	ID              int64              `json:"id"`
	Description     string             `json:"description"`
	Address         ZhkAddress         `json:"address"`
	Header          ZhkHeader          `json:"header"`
	Contacts        ZhkContacts        `json:"contacts"`
	Characteristics ZhkCharacteristics `json:"characteristics"`
	Apartments      ZhkApartments      `json:"apartments"`
	Reviews         ZhkReviews         `json:"reviews"`
}

//easyjson:json
type ZhkAddress struct {
	Address string   `json:"address"`
	Metro   ZhkMetro `json:"metro"`
}

//easyjson:json
type ZhkMetro struct {
	Id      int    `json:"station_id"`
	Color   string `json:"line_color"`
	Line    string `json:"line"`
	Station string `json:"station"`
}

//easyjson:json
type ZhkHeader struct {
	Name         string   `json:"name"`
	LowestPrice  int      `json:"lowest_price"`
	HighestPrice int      `json:"highest_price"`
	Images       []string `json:"images"`
	ImagesSize   int      `json:"images_size"`
}

//easyjson:json
type ZhkContacts struct {
	Developer string `json:"developer"`
	Phone     string `json:"phone"`
}

//easyjson:json
type ZhkCharacteristics struct {
	Decoration    []int        `json:"decoration"`
	Class         string       `json:"class"`
	CeilingHeight CeilingRange `json:"ceiling_height"`
	Floors        FloorsRange  `json:"floors"`
	Square        SquareRange  `json:"square"`
}

//easyjson:json
type CeilingRange struct {
	HighestHeight int `json:"highest_height"`
	LowestHeight  int `json:"lowest_height"`
}

//easyjson:json
type FloorsRange struct {
	HighestFloor int `json:"highest_floor"`
	LowestFloor  int `json:"lowest_floor"`
}

//easyjson:json
type SquareRange struct {
	HighestSquare float64 `json:"highest_square"`
	LowestSquare  float64 `json:"lowest_square"`
}

//easyjson:json
type ZhkApartments struct {
	Apartments []ZhkApartment `json:"items"`
}

//easyjson:json
type ZhkApartment struct {
	HighestPrice int `json:"highest_price"`
	LowestPrice  int `json:"lowest_price"`
	MinSquare    int `json:"min_square"`
	Offers       int `json:"offers"`
	Rooms        int `json:"rooms"`
}

//easyjson:json
type ZhkReviews struct {
	Reviews     []Review `json:"items"`
	Quantity    int      `json:"quantity"`
	TotalRating float64  `json:"total_rating"`
}

//easyjson:json
type Review struct {
	Avatar    string `json:"avatar"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rating    int    `json:"rating"`
	Text      string `json:"tetx"`
}

//easyjson:json
type ZhksInfo []ZhkInfo