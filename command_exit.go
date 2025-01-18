package main

import (
	"os"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
)

func commandExit(cfg *pokeconfig.Config, args ...string) error {
	cfg.Logger.Debug("Exiting Pokedex")
	os.Exit(0)
	return nil
}
