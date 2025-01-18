package commands

import (
	"errors"
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func CommandMapf(cfg *pokedata.Config, args ...string) error {
	cfg.Logger.Debug("Fetching next page of locations")
	locationsResp, err := cfg.PokeapiClient.ListLocations(cfg.NextLocationsURL)
	if err != nil {
		cfg.Logger.Error("Failed to fetch locations: %v", err)
		return err
	}

	cfg.Logger.Info("Found %d locations", len(locationsResp.Results))

	cfg.NextLocationsURL = locationsResp.Next
	cfg.PrevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func CommandMapb(cfg *pokedata.Config, args ...string) error {
	if cfg.PrevLocationsURL == nil {
		cfg.Logger.Error("No previous page available")
		return errors.New("no previous page")
	}

	cfg.Logger.Debug("Fetching previous page of locations")
	locationsResp, err := cfg.PokeapiClient.ListLocations(cfg.PrevLocationsURL)
	if err != nil {
		cfg.Logger.Error("Failed to fetch locations: %v", err)
		return err
	}

	cfg.NextLocationsURL = locationsResp.Next
	cfg.PrevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
