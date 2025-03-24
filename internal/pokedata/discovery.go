package pokedata

// LocationPokemon represents a Pokemon at a specific location
// It's like a coordinate pair that uniquely identifies a discovery

// Stephen: I wonder if there is a cost for how many composite keys will be created? 
type LocationPokemon struct {
	Location string
	Pokemon string
}

// DiscoveryTracker efficiently tracks discovered Pokemon by location
// Think of this as a specialized index system for Pokemon discoveries
type DiscoveryTracker struct {
	// Primary storage - a single flat map using composite keys
	discoveries map[LocationPokemon]bool

	// Indexes for efficient lookups from different perspectives
    // Like having both an author index and a title index in a library
    byLocation map[string]map[string]bool // Location → Pokemon → discovered
    byPokemon  map[string]map[string]bool // Pokemon → Location → discovered
}

// NewDiscoveryTracker creates a new discovery tracker
func NewDiscoveryTracker() *DiscoveryTracker {
	return &DiscoveryTracker{
		discoveries: make(map[LocationPokemon]bool),
		byLocation: make(map[string]map[string]bool),
		byPokemon: make(map[string]map[string]bool),
	}
}

// MarkDiscovered marks a Pokemon as discovered in a location
func (d *DiscoveryTracker) MarkDiscovered(location, pokemon string) {
	// Create the composite key
	key := LocationPokemon{Location: location, Pokemon: pokemon}

	// Skip if already discovered
	if d.discoveries[key] {
		return
	}
	
	// Mark as discovered in our primary storage
	d.discoveries[key] = true

	// Update the location index
	if d.byLocation[location] == nil {
		d.byLocation[location] = make(map[string]bool)
	}
	d.byLocation[location][pokemon] = true

	// Update the Pokemon index
	if d.byPokemon[pokemon] == nil {
		d.byPokemon[pokemon] = make(map[string]bool)
	}
	d.byPokemon[pokemon][location] = true
}

// IsDiscovered checks if a Pokemon has been discovered in a location
func (d *DiscoveryTracker) IsDiscovered(location, pokemon string) bool {
	return d.discoveries[LocationPokemon{Location: location, Pokemon: pokemon}] // TODO I wonder if this is the best way to do this? 
}

// GetDiscoveredLocations returns all locations where a Pokemon has been discovered
func (d *DiscoveryTracker) GetDiscoveredInLocation(location string) []string {
	if d.byLocation[location] == nil {
		return []string{}
	}

	var result []string
	for pokemon := range d.byLocation[location] {
		result = append(result, pokemon)
	}

	return result
}

// GetLocationsWithPokemon returns all locations where a Pokemon was discovered
func (d *DiscoveryTracker) GetLocationsWithPokemon(pokemon string) []string {
	if d.byPokemon[pokemon] == nil {
		return []string{}
	}

	var result []string
	for location := range d.byPokemon[pokemon] {
		result = append(result, location)
	}

	return result
}

// CountDiscoveredInLocation returns the count of discovered Pokemon in a location
func (d *DiscoveryTracker) CountDiscoveredInLocation(location string) int {
	if d.byLocation[location] == nil {
		return 0
	}
	return len(d.byLocation[location])
}

// CountDiscoveredPokemon returns the total count of discovered Pokemon
func (d *DiscoveryTracker) CountDiscoveredPokemon() int {
	if d.byPokemon == nil {
		return 0
	}
	return len(d.byPokemon)
}

// CountTotalDiscoveries returns the total number of discoveries made
func (d *DiscoveryTracker) CountTotalDiscoveries() int {
    return len(d.discoveries)
}



