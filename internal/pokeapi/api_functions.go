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
		c.logger.Debug("Using page URL: %s", *pageURL)
		url = *pageURL
	}

	c.logger.Debug("Fetching locations from: %s", url)

	cacheKey := "location-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for locations")
		locationsResp := PokeAPILocations{}
		err := json.Unmarshal(cachedResp, &locationsResp)
		if err != nil {
			return PokeAPILocations{}, err
		} else {
			return locationsResp, nil
		}

	}

	c.logger.Debug("Cache miss for locations, fetching from API")

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
		c.logger.Error("Area name is required")
		return PokeAPIArea{}, errors.New("area is required")
	}

	cacheKey := "area-pokemon-key-" + url
	c.logger.Debug("Fetching Pokemon for area: %s", area)

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for area: %s", area)
		areaPokemon := PokeAPIArea{}
		err := json.Unmarshal(cachedResp, &areaPokemon)
		if err != nil {
			c.logger.Error("Failed to unmarshal cached area data: %v", err)
			return PokeAPIArea{}, err
		} else {
			return areaPokemon, nil
		}
	}

	c.logger.Debug("Cache miss for area: %s, fetching from API", area)

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

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokeAPIPokemonSpecies{}, err
	}

	return speciesResp, nil
}

func (c *Client) FetchPokemon(pokemonName string) (PokeAPIPokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	// Log the API request attempt
	c.logger.Debug("Attempting to fetch Pokemon: %s", pokemonName)

	if pokemonName == "" {
		return PokeAPIPokemon{}, errors.New("name is required")
	}

	cacheKey := "pokemon-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for Pokemon: %s", pokemonName)
		pokemonResp := PokeAPIPokemon{}
		err := json.Unmarshal(cachedResp, &pokemonResp)
		if err != nil {
			c.logger.Error("Failed to unmarshal cached Pokemon data: %v", err)
			return PokeAPIPokemon{}, err
		} else {
			return pokemonResp, nil
		}
	}

	c.logger.Debug("Cache miss for Pokemon: %s, fetching from API", pokemonName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to fetch Pokemon from API: %v", err)
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

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokeAPIPokemon{}, err
	}

	return pokemonResp, nil
}
