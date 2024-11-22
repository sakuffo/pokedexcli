package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		cfg.logger.Error("Area is required")
		return errors.New("area is required")
	}

	area := args[0]
	cfg.logger.Info("Exploring area: %s", area)

	exploreResp, err := cfg.pokeapiClient.FetchAreaPokemon(area)
	if err != nil {
		cfg.logger.Error("Failed to fetch area %s: %v", area, err)
		return err
	}

	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon: ")
	for _, poke := range exploreResp.PokemonEncounters {
		fmt.Printf("\t- %s\n", poke.Pokemon.Name)
	}

	return nil
}
