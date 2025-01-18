package commands

import (
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func commandPokedex(cfg *pokedata.Config, args ...string) error {
	cfg.Logger.Info("Executing 'pokedex' command")

	fmt.Println("Your Pokedex:")
	fmt.Printf("You have caught %d pokemon\n", len(cfg.CaughtPokemon))
	for _, pokemon := range cfg.CaughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	cfg.Logger.Info("Displayed Pokedex with %d pokemon", len(cfg.CaughtPokemon))
	return nil
}
