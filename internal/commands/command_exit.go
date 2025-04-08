package commands

import (
	"fmt"
	"os"

	"github.com/sakuffo/pokedexcli/internal/config"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func CommandExit(cfg *config.Config, args ...string) error {
	cfg.Logger.Debug("Exiting Pokedex")
	// Save data before exiting
	err := pokedata.SaveData(cfg)
	if err != nil {
		cfg.Logger.Error("Failed to save data: %v", err)
		fmt.Printf("Failed to save data: %v\n", err)
	}
	os.Exit(0)
	return nil
}
