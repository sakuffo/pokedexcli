package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (Locations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		c.logger.Debug("Using page URL: %s", *pageURL)
		url = *pageURL
	}

	c.logger.Debug("Fetching locations from: %s", url)

	cacheKey := "location-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for locations")
		locationsResp := Locations{}
		err := json.Unmarshal(cachedResp, &locationsResp)
		if err != nil {
			return Locations{}, err
		} else {
			return locationsResp, nil
		}

	}

	c.logger.Debug("Cache miss for locations, fetching from API")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Locations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Locations{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Locations{}, err
	}

	locationsResp := Locations{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return Locations{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return Locations{}, err
	}

	return locationsResp, nil
}

func (c *Client) FetchAreaPokemon(area string) (Area, error) {
	url := baseURL + "/location-area/" + area
	if area == "" {
		c.logger.Error("Area name is required")
		return Area{}, errors.New("area is required")
	}

	cacheKey := "area-pokemon-key-" + url
	c.logger.Debug("Fetching Pokemon for area: %s", area)

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for area: %s", area)
		areaPokemon := Area{}
		err := json.Unmarshal(cachedResp, &areaPokemon)
		if err != nil {
			c.logger.Error("Failed to unmarshal cached area data: %v", err)
			return Area{}, err
		} else {
			return areaPokemon, nil
		}
	}

	c.logger.Debug("Cache miss for area: %s, fetching from API", area)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Area{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Area{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Area{}, err
	}

	areaPokemon := Area{}
	err = json.Unmarshal(dat, &areaPokemon)
	if err != nil {
		return Area{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return Area{}, err
	}

	return areaPokemon, nil
}

func (c *Client) FetchPokemonSpecies(name string) (PokemonSpecies, error) {
	url := baseURL + "/pokemon-species/" + name
	if name == "" {
		return PokemonSpecies{}, errors.New("name is required")
	}

	cacheKey := "species-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		speciesResp := PokemonSpecies{}
		err := json.Unmarshal(cachedResp, &speciesResp)
		if err != nil {
			return PokemonSpecies{}, err
		} else {
			return speciesResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonSpecies{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonSpecies{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonSpecies{}, err
	}

	speciesResp := PokemonSpecies{}
	err = json.Unmarshal(dat, &speciesResp)
	if err != nil {
		return PokemonSpecies{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return PokemonSpecies{}, err
	}

	return speciesResp, nil
}

func (c *Client) FetchPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	// Log the API request attempt
	c.logger.Debug("Attempting to fetch Pokemon: %s", pokemonName)

	if pokemonName == "" {
		return Pokemon{}, errors.New("name is required")
	}

	cacheKey := "pokemon-key-" + url

	if cachedResp, ok := c.cache.Get(cacheKey); ok {
		c.logger.Debug("Cache hit for Pokemon: %s", pokemonName)
		pokemonResp := Pokemon{}
		err := json.Unmarshal(cachedResp, &pokemonResp)
		if err != nil {
			c.logger.Error("Failed to unmarshal cached Pokemon data: %v", err)
			return Pokemon{}, err
		} else {
			return pokemonResp, nil
		}
	}

	c.logger.Debug("Cache miss for Pokemon: %s, fetching from API", pokemonName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to fetch Pokemon from API: %v", err)
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonResp := Pokemon{}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		return Pokemon{}, err
	}

	err = c.cache.Add(cacheKey, dat)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemonResp, nil
}
