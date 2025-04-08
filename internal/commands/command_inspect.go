package commands

import (
	"errors"
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/config"
)

func CommandInspect(cfg *config.Config, args ...string) error {
	if len(args) != 1 {
		cfg.Logger.Error("Inspect command called without a pokemon name")
		return errors.New("inspect requires a pokemon name")
	}

	pokemonName := args[0]
	cfg.Logger.Info("Executing 'inspect' command for %s", pokemonName)

	pokemon, exists := cfg.CaughtPokemon[pokemonName]
	if !exists {
		cfg.Logger.Error("Attempted to inspect uncaught Pokemon: %s", pokemonName)
		return errors.New("you haven't caught this pokemon yet")
	}

	cfg.Logger.Info("Displaying details for Pokemon: %s", pokemonName)

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: \033[32m%d\033[0m\n", pokemon.Height)
	fmt.Printf("Weight: \033[32m%d\033[0m\n", pokemon.Weight)
	fmt.Printf("Species: \033[32m%s\033[0m\n", pokemon.Species)
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
