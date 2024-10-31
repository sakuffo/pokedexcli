package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("inspect requires a pokemon name")
	}

	if args[0] == "all" {
		for _, pokemon := range cfg.pokedex.Pokemon {
			fmt.Printf("%s\n", pokemon.Name)
		}

		return nil
	}

	pokemonName := args[0]

	pokemon, exists := cfg.pokedex.Pokemon[pokemonName]
	if !exists {
		return errors.New("pokemon not found")
	}

	fmt.Printf("%s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Types: %v\n", pokemon.Types[0].Type.Name)
	fmt.Printf("Species: %v\n", pokemon.Species.Name)

	return nil
}
