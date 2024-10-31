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

func (c *Client) FetchPokemonSpecies(name string) (PokeAPIPokemonSpecies, error) {
	url := baseURL + "/pokemon-species/" + name
	if name == "" {
		return PokeAPIPokemonSpecies{}, errors.New("name is required")
	}

	cacheKey := "species-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		speciesResp := PokeAPIPokemonSpecies{}
		err := json.Unmarshal(cachedResp, &speciesResp)
		if err != nil {
			return PokeAPIPokemonSpecies{}, err
		} else {
			return speciesResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeAPIPokemonSpecies{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeAPIPokemonSpecies{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeAPIPokemonSpecies{}, err
	}

	speciesResp := PokeAPIPokemonSpecies{}
	err = json.Unmarshal(dat, &speciesResp)
	if err != nil {
		return PokeAPIPokemonSpecies{}, err
	}

	return speciesResp, nil
}

func (c *Client) FetchPokemon(name string) (PokeAPIPokemon, error) {
	url := baseURL + "/pokemon/" + name
	if name == "" {
		return PokeAPIPokemon{}, errors.New("name is required")
	}

	cacheKey := "pokemon-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		pokemonResp := PokeAPIPokemon{}
		err := json.Unmarshal(cachedResp, &pokemonResp)
		if err != nil {
			return PokeAPIPokemon{}, err
		} else {
			return pokemonResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeAPIPokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeAPIPokemon{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeAPIPokemon{}, err
	}

	pokemonResp := PokeAPIPokemon{}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		return PokeAPIPokemon{}, err
	}

	return pokemonResp, nil
}
