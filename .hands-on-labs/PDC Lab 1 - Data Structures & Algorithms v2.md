# Golang Senior Engineer Lab: Data Structures & Nested Maps

## Scenario

You've been tasked with reviewing the code of an intermediate Go engineer working on the PokédexCLI project. One particular area of concern is how the application tracks discovered Pokémon across different locations. The current implementation uses nested maps, which works but has several inefficiencies. Your job is to coach the developer on a more idiomatic and efficient approach.

## Current Implementation Review

Let's examine the current code in `command_explore.go`:

```go
// Initialize the area in discovered map if it doesn't exist
if _, exists := cfg.DiscoveredPokemon[locationName]; !exists {
    cfg.DiscoveredPokemon[locationName] = make(map[string]bool)
}

// Get all undiscovered pokemon in this area
var undiscovered []string
discoveredInArea := 0
totalInArea := len(exploreResp.PokemonEncounters)

for _, poke := range exploreResp.PokemonEncounters {
    if !cfg.DiscoveredPokemon[locationName][poke.Pokemon.Name] {
        undiscovered = append(undiscovered, poke.Pokemon.Name)
    } else {
        discoveredInArea++
    }
}
```

And in `pokeconfig.go`:

```go
type Config struct {
    // ...other fields...
    DiscoveredPokemon map[string]map[string]bool
    // ...other fields...
}
```

## Conceptual Understanding: Nested Maps vs. Domain Types

### The Problem with Nested Maps

Think of our current data structure like a disorganized filing cabinet. For each location (the outer drawer), we have another set of folders (the inner map), and in each folder, we just have a simple "yes" marker for each Pokémon. This approach presents several challenges:

1. **Memory Fragmentation**: Each inner map requires its own memory allocation, potentially leading to memory fragmentation, especially as the application scales.
    
2. **Defensive Coding Overhead**: We constantly need to check if the inner map exists before using it, adding repetitive boilerplate code.
    
3. **Limited Querying Capabilities**: Answering questions like "Where have I discovered Pikachu?" requires iterating through every location.
    
4. **Data Integrity Risks**: With direct access to the nested maps, it's easy for code in various places to manipulate the structure inconsistently.
    
5. **Serialization Complexity**: When saving this structure to disk, nested maps may not serialize/deserialize as efficiently as a flat structure.
    

It's like having a spreadsheet where you need to first check if a row exists before you can access its cells. This creates unnecessary complexity throughout your codebase.

### A Better Approach: Domain-Specific Types

In Go, we typically favor creating dedicated types that represent our domain concepts rather than using generic data structures directly. This is similar to how a database uses indexes to optimize different query patterns.

Think of it like designing a specialized card catalog system for a library instead of just piling books on shelves. The specialized system allows you to look up books by author, title, or subject without searching through everything.

## Exercise: Implementing a Discovery Tracker

Let's create a domain-specific type that efficiently tracks Pokémon discoveries:

### Step 1: Create a dedicated discovery tracker

Create a new file `internal/pokedata/discovery.go`:

```go
package pokedata

// LocationPokemon represents a Pokemon at a specific location
// It's like a coordinate pair that uniquely identifies a discovery
type LocationPokemon struct {
    Location string
    Pokemon  string
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
        byLocation:  make(map[string]map[string]bool),
        byPokemon:   make(map[string]map[string]bool),
    }
}

// MarkDiscovered marks a Pokemon as discovered in a location
func (d *DiscoveryTracker) MarkDiscovered(location, pokemon string) {
    // Create our composite key
    key := LocationPokemon{Location: location, Pokemon: pokemon}
    
    // Skip if already discovered
    if d.discoveries[key] {
        return
    }
    
    // Mark as discovered in our primary storage
    d.discoveries[key] = true
    
    // Update location index
    if d.byLocation[location] == nil {
        d.byLocation[location] = make(map[string]bool)
    }
    d.byLocation[location][pokemon] = true
    
    // Update Pokemon index
    if d.byPokemon[pokemon] == nil {
        d.byPokemon[pokemon] = make(map[string]bool)
    }
    d.byPokemon[pokemon][location] = true
}

// IsDiscovered checks if a Pokemon is discovered in a location
func (d *DiscoveryTracker) IsDiscovered(location, pokemon string) bool {
    return d.discoveries[LocationPokemon{Location: location, Pokemon: pokemon}]
}

// GetDiscoveredInLocation returns all discovered Pokemon in a location
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

// CountTotalDiscoveries returns the total number of discoveries made
func (d *DiscoveryTracker) CountTotalDiscoveries() int {
    return len(d.discoveries)
}
```

### Step 2: Update the Config struct to use our new tracker

Modify `internal/pokedata/pokeconfig.go`:

```go
type Config struct {
    // ... existing fields
    
    // Replace this:
    // DiscoveredPokemon map[string]map[string]bool
    
    // With this:
    Discoveries *DiscoveryTracker
}

func New(logLevel logger.LogLevel) *Config {
    // ... existing code
    
    // Initialize our discovery tracker
    discoveries := NewDiscoveryTracker()
    
    // Migrate existing data if it exists
    if data.DiscoveredPokemon != nil {
        for location, pokemons := range data.DiscoveredPokemon {
            for pokemon, discovered := range pokemons {
                if discovered {
                    discoveries.MarkDiscovered(location, pokemon)
                }
            }
        }
    }
    
    cfg.Discoveries = discoveries
    
    // ... rest of initialization
    return cfg
}
```

### Step 3: Update the explore command

Modify `internal/commands/command_explore.go`:

```go
func CommandExplore(cfg *pokedata.Config, args ...string) error {
    // ... existing code
    
    // Replace the old discovery logic:
    var undiscovered []string
    totalInArea := len(exploreResp.PokemonEncounters)
    
    for _, poke := range exploreResp.PokemonEncounters {
        if !cfg.Discoveries.IsDiscovered(locationName, poke.Pokemon.Name) {
            undiscovered = append(undiscovered, poke.Pokemon.Name)
        }
    }
    
    discoveredInArea := cfg.Discoveries.CountDiscoveredInLocation(locationName)
    
    // ... existing exploration code ...
    
    // When marking newly discovered Pokemon:
    for i := 0; i < numToDiscover; i++ {
        pokeName := undiscovered[i]
        cfg.Discoveries.MarkDiscovered(locationName, pokeName)
        newlyDiscovered[pokeName] = true
        discoveredInArea++
    }
    
    // ... rest of the code
}
```

### Step 4: Update data persistence

You'll also need to adjust how discovery data is saved. In `internal/pokedata/data_funcs.go`:

```go
// Add a method to convert DiscoveryTracker to the old format for backward compatibility
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
    // ... existing code
    
    data := &Data{
        CaughtPokemon:     cfg.CaughtPokemon,
        PartyMembers:      cfg.Party.Members,
        DiscoveredPokemon: cfg.Discoveries.ToMap(),  // Convert to map for storage
    }
    
    err := cfg.Persistence.Save(data)
    // ... rest of the code
}
```

## Understanding the Benefits

This refactoring provides multiple benefits:

1. **Encapsulation**: We've hidden the internal structure behind a clean API, making it harder to introduce bugs.
    
2. **Efficiency**: We only create inner maps when needed, reducing memory overhead.
    
3. **New Query Capabilities**: We can now easily answer questions like "Where have I seen Pikachu?" without looping through all locations.
    
4. **Semantic Clarity**: The code now expresses intent through descriptive method names rather than raw data manipulation.
    
5. **Single Responsibility**: The discovery tracking logic is centralized in one place rather than scattered throughout commands.
    

Think of this approach like upgrading from a pile of sticky notes to a digital contact manager. The contact manager doesn't just store information – it organizes it in ways that make specific queries fast and easy.

## Key Go Patterns Demonstrated

1. **Domain-Driven Design**: Creating types that model the domain concepts rather than using generic data structures.
    
2. **Composition**: Building complex functionality through simple, composed types.
    
3. **Information Hiding**: Exposing only what's needed through method interfaces.
    
4. **Lazy Initialization**: Only creating resources (inner maps) when they're actually needed.
    
5. **Indexed Data**: Maintaining multiple views of the same data for efficient access patterns.
    

These patterns are cornerstone practices for writing maintainable, production-grade Go code and will serve you well as your projects grow in complexity.