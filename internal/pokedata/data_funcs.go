package pokedata

import (
	"errors"
)

// Add a method to convert DiscoveryTracker to a map for storage in JSON
func (d *DiscoveryTracker) ToMap() map[string]map[string]bool {
	result := make(map[string]map[string]bool)

	for location, pokemons := range d.byLocation {
		result[location] = make(map[string]bool)
		for pokemon := range pokemons {
			result[location][pokemon] = true
		}
	}

	return result
}

func SaveData(cfg *Config) error {
	cfg.Logger.Debug("Saving data")
	if cfg.Persistence == nil {
		cfg.Logger.Error("Persistence not initialized")
		return errors.New("persistence not initialized")
	}

	data := &Data{
		CaughtPokemon:   cfg.CaughtPokemon,
		PartyMembers:    cfg.Party.Members,
		DiscoveriesMade: cfg.Discoveries.ToMap(), // Convert to map for JSON serialization & storage
	}

	err := cfg.Persistence.Save(data)
	if err != nil {
		cfg.Logger.Error("Failed to save data: %v", err)
		return err
	}

	return nil
}
