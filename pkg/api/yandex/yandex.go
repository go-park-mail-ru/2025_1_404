package yandex

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type yandexRepo struct {
	apiKey string
}

func New(cfg *config.YandexConfig) *yandexRepo {
	return &yandexRepo{apiKey: cfg.Token}
}

func (repo *yandexRepo) GetCoordinatesOfAddress(address string) (*Coordinates, error) {
	geocodeResponse, err := repo.GeocodeAddress(address)
	if err != nil {
		return nil, err
	}
	if len(geocodeResponse.GeoObjectCollection.FeatureMembers) == 0 {
		return nil, fmt.Errorf("no coordinates found for address: %s", address)
	}
	point := geocodeResponse.GeoObjectCollection.FeatureMembers[0].GeoObject.Point.Pos
	coords := strings.Split(point, " ")
	if len(coords) != 2 {
		return nil, fmt.Errorf("invalid coordinates format: %s", point)
	}
	latitude, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude value: %s", coords[1])
	}
	longitude, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude value: %s", coords[0])
	}
	return &Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func (repo *yandexRepo) GeocodeAddress(address string) (*GeocodeResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://geocode-maps.yandex.ru/v1", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("apikey", repo.apiKey)
	q.Add("geocode", address)
	q.Add("format", "json")
	q.Add("results", "1")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get geocode: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	apiResponse := &ApiResponse{}
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse.Response, nil
}
