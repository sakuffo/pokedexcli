package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		cfg.logger.Error("Inspect command called without a pokemon name")
		return errors.New("inspect requires a pokemon name")
	}

	pokemonName := args[0]
	cfg.logger.Info("Executing 'inspect' command for %s", pokemonName)

	pokemon, exists := cfg.caughtPokemon[pokemonName]
	if !exists {
		cfg.logger.Error("Attempted to inspect uncaught Pokemon: %s", pokemonName)
		return errors.New("You haven't caught this pokemon yet")
	}

	cfg.logger.Info("Displaying details for Pokemon: %s", pokemonName)

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: \033[32m%d\033[0m\n", pokemon.Height)
	fmt.Printf("Weight: \033[32m%d\033[0m\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: \033[32m%d\033[0m\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	fmt.Println()

	return nil
}
