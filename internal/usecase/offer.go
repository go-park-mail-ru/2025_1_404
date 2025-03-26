package usecase

import (
	"github.com/go-park-mail-ru/2025_1_404/domain"
	"github.com/go-park-mail-ru/2025_1_404/internal/mocks"
)

// GetOffers Получение списка объявлений (пока что из моков)
func GetOffers() []domain.Offer {
	return mocks.GetMockOffers()
}
