package main

import "fmt"

func commandPokedex(cfg *config, args ...string) error {
	cfg.logger.Info("Executing 'pokedex' command")

	fmt.Println("Your Pokedex:")
	fmt.Printf("You have caught %d pokemon\n", len(cfg.caughtPokemon))
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	cfg.logger.Info("Displayed Pokedex with %d pokemon", len(cfg.caughtPokemon))
	return nil
}
