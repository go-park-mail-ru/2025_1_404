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

// ZhkInfo Структура с информацией о ЖК
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

// ZhkAddress Расположение ЖК
type ZhkAddress struct {
	Address string   `json:"address"`
	Metro   ZhkMetro `json:"metro"`
}

// ZhkMetro Метро ЖК
type ZhkMetro struct {
	Id      int    `json:"station_id"`
	Color   string `json:"line_color"`
	Line    string `json:"line"`
	Station string `json:"station"`
}

// ZhkHeader Заголовок ЖК
type ZhkHeader struct {
	Name         string   `json:"name"`
	LowestPrice  int      `json:"lowest_price"`
	HighestPrice int      `json:"highest_price"`
	Images       []string `json:"images"`
	ImagesSize   int      `json:"images_size"`
}

// ZhkContacts Контактная информация
type ZhkContacts struct {
	Developer string `json:"developer"`
	Phone     string `json:"phone"`
}

// ZhkCharacteristics Характеристики ЖК
type ZhkCharacteristics struct {
	Decoration    []int        `json:"decoration"`
	Class         string       `json:"class"`
	CeilingHeight CeilingRange `json:"ceiling_height"`
	Floors        FloorsRange  `json:"floors"`
	Square        SquareRange  `json:"square"`
}

// CeilingRange Диапазон потолков
type CeilingRange struct {
	HighestHeight int `json:"highest_height"`
	LowestHeight  int `json:"lowest_height"`
}

// FloorsRange Диапазон этажей
type FloorsRange struct {
	HighestFloor int `json:"highest_floor"`
	LowestFloor  int `json:"lowest_floor"`
}

// SquareRange Диапазон площадей
type SquareRange struct {
	HighestSquare float64 `json:"highest_square"`
	LowestSquare  float64 `json:"lowest_square"`
}

// ZhkApartments Предложения ЖК
type ZhkApartments struct {
	Apartments []ZhkApartment `json:"items"`
}

// ZhkApartment Предложение ЖК
type ZhkApartment struct {
	HighestPrice int `json:"highest_price"`
	LowestPrice  int `json:"lowest_price"`
	MinSquare    int `json:"min_square"`
	Offers       int `json:"offers"`
	Rooms        int `json:"rooms"`
}

// ZhkReview Отзывы на ЖК
type ZhkReviews struct {
	Reviews     []Review `json:"items"`
	Quantity    int      `json:"quantity"`
	TotalRating float64  `json:"total_rating"`
}

// Review Отзыв на ЖК
type Review struct {
	Avatar    string `json:"avatar"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rating    int    `json:"rating"`
	Text      string `json:"tetx"`
}
