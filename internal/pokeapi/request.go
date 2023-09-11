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
	PokemonEndpoint      = "https://pokeapi.co/api/v2/pokemon/"
)

type PokeJSON interface {
	LocationData | LocationAreaData | PokemonData
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

func parse[T PokeJSON](responseBody []byte, data *T) error {
	if err := json.Unmarshal(responseBody, data); err != nil {
		return err
	}
	return nil
}

func parseAndReturn[T PokeJSON](body []byte, data T) (T, error) {
	err := parse(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
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

func Get[T PokeJSON](url string, data T, cache *pokecache.Cache) (T, error) {
	body, ok := cache.Get(url)
	if ok {
		return parseAndReturn(body, data)
	}
	body, err := fetchAndRead(url, cache)
	if err != nil {
		return data, err
	}
	return parseAndReturn(body, data)
}
