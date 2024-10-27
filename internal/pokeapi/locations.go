package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (PokeAPILocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	cacheKey := "location-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		locationsResp := PokeAPILocations{}
		err := json.Unmarshal(cachedResp, &locationsResp)
		if err != nil {
			return PokeAPILocations{}, err
		} else {
			return locationsResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeAPILocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeAPILocations{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeAPILocations{}, err
	}

	locationsResp := PokeAPILocations{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return PokeAPILocations{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokeAPILocations{}, err
	}

	return locationsResp, nil
}

func (c *Client) FetchAreaPokemon(area string) (PokeAPIArea, error) {
	url := baseURL + "/location-area/" + area
	if area == "" {
		return PokeAPIArea{}, errors.New("area is required")
	}

	cacheKey := "pokemon-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		areaPokemon := PokeAPIArea{}
		err := json.Unmarshal(cachedResp, &areaPokemon)
		if err != nil {
			return PokeAPIArea{}, err
		} else {
			return areaPokemon, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeAPIArea{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeAPIArea{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeAPIArea{}, err
	}

	areaPokemon := PokeAPIArea{}
	err = json.Unmarshal(dat, &areaPokemon)
	if err != nil {
		return PokeAPIArea{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokeAPIArea{}, err
	}

	return areaPokemon, nil
}
