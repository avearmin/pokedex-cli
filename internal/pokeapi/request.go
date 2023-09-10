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

type PokeJSON interface {
	LocationData | LocationAreaData
}

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

type LocationAreaData struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	Location             Location               `json:"location"`
	Names                []Names                `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterMethodRatesVersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}
type EncounterMethodRates struct {
	EncounterMethod EncounterMethod                      `json:"encounter_method"`
	VersionDetails  []EncounterMethodRatesVersionDetails `json:"version_details"`
}
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Names struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}
type PokemonEncountersVersionDetails struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"max_chance"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon                           `json:"pokemon"`
	VersionDetails []PokemonEncountersVersionDetails `json:"version_details"`
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
