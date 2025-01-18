package pokedata

import (
	"errors"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
)

func SaveData(cfg *pokeconfig.Config) error {
	cfg.Logger.Debug("Saving data")
	if cfg.Persistence == nil {
		cfg.Logger.Error("Persistence not initialized")
		return errors.New("persistence not initialized")
	}

	data := &Data{
		CaughtPokemon: cfg.CaughtPokemon,
		PartyMembers:  cfg.Party.Members,
	}

	err := cfg.Persistence.Save(data)
	if err != nil {
		cfg.Logger.Error("Failed to save data: %v", err)
		return err
	}

	return nil
}
