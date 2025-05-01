package yandex

type YandexRepo interface {
	GeocodeAddress(address string) (*GeocodeResponse, error)
	GetCoordinatesOfAddress(address string) (*Coordinates, error)
}
