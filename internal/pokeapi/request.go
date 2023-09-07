package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	pokeAPIUrl   = "https://pokeapi.co/api/v2/"
	LocationArea = "location-area/"
)

type LocationData struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []LocationResult `json:"results"`
}

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func sendGETRequest(pathComponent string) (*http.Response, error) {
	endpoint := pokeAPIUrl + pathComponent
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Response failed with status code: %d", res.StatusCode)
	}
	return res, nil
}

func readResponseBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parse(responseBody []byte, locationData *LocationData) error {
	if err := json.Unmarshal(responseBody, locationData); err != nil {
		return err
	}
	return nil
}

func Get(pathComponent string) (LocationData, error) {
	var locationData LocationData
	res, err := sendGETRequest(pathComponent)
	if err != nil {
		return locationData, err
	}
	body, err := readResponseBody(res)
	if err != nil {
		return locationData, err
	}
	err = parse(body, &locationData)
	if err != nil {
		return locationData, err
	}
	return locationData, nil
}
