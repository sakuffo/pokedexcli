package pokeapi

import (
	"encoding/json"
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
