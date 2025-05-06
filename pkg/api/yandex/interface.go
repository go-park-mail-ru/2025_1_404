package yandex

//go:generate mockgen -source interface.go -destination=mocks/interface.go -package=mocks

type YandexRepo interface {
	GeocodeAddress(address string) (*GeocodeResponse, error)
	GetCoordinatesOfAddress(address string) (*Coordinates, error)
}
