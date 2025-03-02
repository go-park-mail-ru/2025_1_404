package main

// User представляет зарегистрированного пользователя
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"` // Храним хешированный пароль
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsRealtor bool   `json:"is_realtor"`
}

// Offer представляет собой объявление о квартире
type Offer struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Address    string   `json:"address"`
	Area       int      `json:"area"`       // Площадь в м²
	Rooms      int      `json:"rooms"`      // Количество комнат
	Floor      int      `json:"floor"`      // Этаж
	Renovation string   `json:"renovation"` // Тип ремонта
	DealType   string   `json:"deal_type"`  // Тип сделки (продажа, аренда)
	HouseType  string   `json:"house_type"` // Тип жилья (новостройка, вторичка)
	AdType     string   `json:"ad_type"`    // Тип объявления (агентство, собственник)
	Seller     User     `json:"seller"`     // Продавец (ссылка на пользователя)
	Images     []string `json:"images"`     // Список изображений
}
