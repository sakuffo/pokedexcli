package discovery

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// LocationPokemon represents a unique Pokemon species found at a specific location.
// Using a struct ensures type safety compared to simple strings.
type LocationPokemon struct {
	LocationName string
	PokemonName  string
}

// DiscoveryTracker keeps track of which Pokemon have been discovered in which locations.
// It uses a map with LocationPokemon as the key for type safety and clarity.
type DiscoveryTracker struct {
	// discovered uses LocationPokemon struct as key for type safety.
	// The bool value is always true; presence in the map indicates discovery.
	discovered map[LocationPokemon]bool
	mu         sync.RWMutex
}

// NewDiscoveryTracker initializes a new DiscoveryTracker.
func NewDiscoveryTracker() *DiscoveryTracker {
	return &DiscoveryTracker{
		discovered: make(map[LocationPokemon]bool),
	}
}

// MarkDiscovered records that a specific Pokemon has been discovered at a location.
func (dt *DiscoveryTracker) MarkDiscovered(locationName, pokemonName string) {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	key := LocationPokemon{LocationName: locationName, PokemonName: pokemonName}
	dt.discovered[key] = true
}

// IsDiscovered checks if a specific Pokemon has been discovered at a location.
func (dt *DiscoveryTracker) IsDiscovered(locationName, pokemonName string) bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	key := LocationPokemon{LocationName: locationName, PokemonName: pokemonName}
	_, found := dt.discovered[key]
	return found
}

// GetProgress returns the number of unique Pokemon discovered in a given location
// and the total number of unique Pokemon species encountered across all locations so far.
func (dt *DiscoveryTracker) GetProgress(locationName string) (discoveredInLocation int, totalUniqueSpecies int) {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	locationDiscoveries := make(map[string]bool) // Pokemon names in this location
	totalSpecies := make(map[string]bool)        // Unique pokemon names across all locations

	for key := range dt.discovered {
		totalSpecies[key.PokemonName] = true
		if key.LocationName == locationName {
			locationDiscoveries[key.PokemonName] = true
		}
	}

	return len(locationDiscoveries), len(totalSpecies)
}

// ListDiscoveriesForLocation returns a sorted list of discovered Pokemon names for a specific location.
func (dt *DiscoveryTracker) ListDiscoveriesForLocation(locationName string) []string {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	var names []string
	for key := range dt.discovered {
		if key.LocationName == locationName {
			names = append(names, key.PokemonName)
		}
	}
	sort.Strings(names) // Ensure consistent order
	return names
}

// ToMap converts the internal discovery map (using struct keys) to a format
// suitable for JSON serialization (map[string]map[string]bool).
func (dt *DiscoveryTracker) ToMap() map[string]map[string]bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	serializableMap := make(map[string]map[string]bool)
	for key := range dt.discovered {
		if _, ok := serializableMap[key.LocationName]; !ok {
			serializableMap[key.LocationName] = make(map[string]bool)
		}
		serializableMap[key.LocationName][key.PokemonName] = true
	}
	return serializableMap
}

// GetDiscoveredInLocation returns all discovered Pokemon in a location.
// This is an alias for ListDiscoveriesForLocation for backward compatibility.
func (dt *DiscoveryTracker) GetDiscoveredInLocation(location string) []string {
	return dt.ListDiscoveriesForLocation(location)
}

// CountDiscoveredInLocation returns the count of discovered Pokemon in a location
func (dt *DiscoveryTracker) CountDiscoveredInLocation(location string) int {
	count, _ := dt.GetProgress(location)
	return count
}

// String provides a simple string representation, mainly for debugging.
func (dt *DiscoveryTracker) String() string {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return fmt.Sprintf("DiscoveryTracker with %d entries", len(dt.discovered))
}

// NOTE: FromMap is not strictly needed if Load handles migration,
// but could be useful for testing or other scenarios.
// func FromMap(data map[string]map[string]bool) *DiscoveryTracker {
// 	dt := NewDiscoveryTracker()
// 	dt.mu.Lock() // Lock for initial population
// 	defer dt.mu.Unlock()
// 	for location, pokemons := range data {
// 		for pokemon := range pokemons {
// 			key := LocationPokemon{LocationName: location, PokemonName: pokemon}
// 			dt.discovered[key] = true
// 		}
// 	}
// 	return dt
// }

// IsInitialized checks if the discovery tracker's map has been initialized.
func (dt *DiscoveryTracker) IsInitialized() bool {
	return dt.discovered != nil
}

// MarshalJSON implements json.Marshaler interface for proper serialization.
func (dt *DiscoveryTracker) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler interface for proper deserialization.
func (dt *DiscoveryTracker) UnmarshalJSON(data []byte) error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	// Initialize the map if needed
	if dt.discovered == nil {
		dt.discovered = make(map[LocationPokemon]bool)
	}

	// Parse the JSON into a temporary map
	var tempMap map[string]map[string]bool
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// Convert the map back to our complex key structure
	for location, pokemons := range tempMap {
		for pokemon := range pokemons {
			key := LocationPokemon{LocationName: location, PokemonName: pokemon}
			dt.discovered[key] = true
		}
	}

	return nil
}
