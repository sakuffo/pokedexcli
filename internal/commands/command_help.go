package commands

import (
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/config"
)

func CommandHelp(cfg *config.Config, args ...string) error {
	cfg.Logger.Debug("Executing 'help' command")
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range GetCommands() {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}

	fmt.Println()
	cfg.Logger.Debug("Displayed help")
	return nil
}
