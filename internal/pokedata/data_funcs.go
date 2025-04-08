package pokedata

import (
	"errors"

	"github.com/sakuffo/pokedexcli/internal/config"
	"github.com/sakuffo/pokedexcli/internal/persistence"
)

// SaveData prepares the data from the config and saves it using the Persistence service.
// WARNING: This function might be moved or redesigned during later refactoring steps.
func SaveData(cfg *config.Config) error {
	cfg.Logger.Debug("Preparing to save application data...")
	if cfg.Persistence == nil {
		cfg.Logger.Error("SaveData failed: Persistence service is not initialized in config")
		return errors.New("persistence service not initialized")
	}

	// Create the data structure to be saved.
	dataToSave := &persistence.Data{
		CaughtPokemon: cfg.CaughtPokemon,
		PartyMembers:  cfg.Party.Members,
		Discoveries:   cfg.Discoveries,
	}

	// Call the Save method on the persistence service
	err := cfg.Persistence.Save(dataToSave)
	if err != nil {
		// Error already logged by Persistence.Save, just return it
		return err // Return the original error for the caller (e.g., REPL)
	}

	cfg.Logger.Info("Application data saved successfully.")
	return nil
}
