package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/avearmin/pokedex-cli/internal/pokecache"
)

const (
	LocationAreaEndpoint = "https://pokeapi.co/api/v2/location-area/"
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

func fetch(url string) (*http.Response, error) {
	res, err := http.Get(url)
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

func fetchAndRead(url string, cache *pokecache.Cache) ([]byte, error) {
	res, err := fetch(url)
	if err != nil {
		return nil, err
	}
	body, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	cache.Add(url, body)
	return body, nil
}

func Get(url string, cache *pokecache.Cache) (LocationData, error) {
	var locationData LocationData
	body, ok := cache.Get(url)
	if ok {
		err := parse(body, &locationData)
		if err != nil {
			return locationData, err
		}
		return locationData, nil
	}
	body, err := fetchAndRead(url, cache)
	if err != nil {
		return locationData, err
	}
	err = parse(body, &locationData)
	if err != nil {
		return locationData, err
	}
	return locationData, nil
}
