package main

import "net/http"

var offers = []Offer{
	{
		ID: 1, Title: "Квартира у метро", Address: "ул. Ленина, 10", Area: 50, Rooms: 2, Floor: 3,
		Renovation: "Евроремонт", DealType: "Продажа", HouseType: "Вторичка", AdType: "Собственник",
		Seller: User{ID: 1, Email: "default@example.com", FirstName: "Иван", LastName: "Иванов", IsRealtor: false},
		Images: []string{
			"https://img.dmclk.ru/vitrina/owner/a8/28/a8285bf0f4e94f0ba3c894cc40eb38c5.webp",
			"https://img.dmclk.ru/vitrina/owner/a8/d2/a8d2cefde0fd4838b7f2adbec33adc6b.webp",
		},
	},
}

// Получение списка объявлений
func getOffers(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, offers, http.StatusOK)
}
