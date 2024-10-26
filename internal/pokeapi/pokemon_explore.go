package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (c *Client) FetchAreaPokemon(area string) (PokeAPIAreaPokemon, error) {
	url := baseURL + "/location-area/" + area
	if area == "" {
		return PokeAPIAreaPokemon{}, errors.New("area is required")
	}

	cacheKey := "pokemon-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		areaPokemon := PokeAPIAreaPokemon{}
		err := json.Unmarshal(cachedResp, &areaPokemon)
		if err != nil {
			return PokeAPIAreaPokemon{}, err
		} else {
			return areaPokemon, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeAPIAreaPokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeAPIAreaPokemon{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeAPIAreaPokemon{}, err
	}

	areaPokemon := PokeAPIAreaPokemon{}
	err = json.Unmarshal(dat, &areaPokemon)
	if err != nil {
		return PokeAPIAreaPokemon{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokeAPIAreaPokemon{}, err
	}

	return areaPokemon, nil
}
