package yandex

type Point struct {
	Pos string `json:"pos"`
}

type GeoObject struct {
	Point Point `json:"Point"`
}

type FeatureMember struct {
	GeoObject GeoObject `json:"GeoObject"`
}

type GeoObjectCollection struct {
	FeatureMembers []FeatureMember `json:"featureMember"`
}

type GeocodeResponse struct {
	GeoObjectCollection GeoObjectCollection `json:"GeoObjectCollection"`
}

type ApiResponse struct {
	Response GeocodeResponse `json:"response"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
