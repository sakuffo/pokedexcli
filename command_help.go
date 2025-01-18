package main

import (
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
)

func commandHelp(cfg *pokeconfig.Config, args ...string) error {
	cfg.Logger.Debug("Executing 'help' command")
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	cfg.Logger.Debug("Displayed help")
	return nil
}
