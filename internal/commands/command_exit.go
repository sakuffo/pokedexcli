package commands

import (
	"os"

	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func commandExit(cfg *pokedata.Config, args ...string) error {
	cfg.Logger.Debug("Exiting Pokedex")
	os.Exit(0)
	return nil
}
