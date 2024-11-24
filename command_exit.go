package main

import "os"

func commandExit(cfg *config, args ...string) error {
	cfg.logger.Debug("Exiting Pokedex")
	os.Exit(0)
	return nil
}
