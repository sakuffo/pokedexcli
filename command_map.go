package main

import (
	"errors"
	"fmt"
)

func commandMapf(cfg *config, args ...string) error {
	cfg.logger.Debug("Fetching next page of locations")
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		cfg.logger.Error("Failed to fetch locations: %v", err)
		return err
	}

	cfg.logger.Info("Found %d locations", len(locationsResp.Results))

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		cfg.logger.Error("No previous page available")
		return errors.New("no previous page")
	}

	cfg.logger.Debug("Fetching previous page of locations")
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		cfg.logger.Error("Failed to fetch locations: %v", err)
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
