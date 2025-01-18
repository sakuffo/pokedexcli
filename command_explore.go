package main

import (
	"errors"
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
)

func commandExplore(cfg *pokeconfig.Config, args ...string) error {

	if len(args) != 1 {
		cfg.Logger.Error("Area is required")
		return errors.New("area is required")
	}

	area := args[0]
	cfg.Logger.Info("Exploring area: %s", area)

	exploreResp, err := cfg.PokeapiClient.FetchAreaPokemon(area)
	if err != nil {
		cfg.Logger.Error("Failed to fetch area %s: %v", area, err)
		return err
	}

	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon: ")
	for _, poke := range exploreResp.PokemonEncounters {
		fmt.Printf("\t- %s\n", poke.Pokemon.Name)
	}

	return nil
}
