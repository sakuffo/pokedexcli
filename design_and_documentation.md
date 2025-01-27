# PokédexCLI Design Documentation
as of 12/31/2024

## Overview
PokédexCLI is a command-line interface application that simulates a Pokémon trainer's journey. It allows users to explore locations, catch Pokémon, manage their party, and interact with their Pokédex. The application is built in Go and integrates with the PokéAPI for Pokémon data.

## Core Architecture

### Main Components

1. **REPL (Read-Eval-Print Loop)**
   - Entry point of the application
   - Handles user input and command routing
   - Manages application state through the `config` struct
   - Implements graceful shutdown with signal handling

2. **Configuration Management**
   - Central `config` struct maintains application state
   - Manages:
     - PokéAPI client
     - Location navigation
     - Caught Pokémon
     - Party system
     - Data persistence
     - Logging

3. **Command System**
   - Modular command structure with individual command files
   - Each command implements a callback function
   - Available commands:
     - `help`: Display available commands
     - `exit`: Exit the application
     - `map`: Navigate location areas
     - `explore`: List Pokémon in an area
     - `catch`: Attempt to catch a Pokémon
     - `inspect`: View Pokémon details
     - `pokedex`: List caught Pokémon
     - `party`: Manage Pokémon party

### Internal Packages

1. **pokeapi**
   - Handles communication with the PokéAPI
   - Implements rate limiting and timeout management
   - Defines type structures for API responses
   - Provides methods for fetching:
     - Pokémon data
     - Location areas
     - Species information

2. **pokecache**
   - Implements in-memory caching system
   - Features:
     - Automatic cache cleanup
     - Thread-safe operations
     - Configurable cleanup intervals
     - Cache hit/miss logging

3. **pokedata**
   - Manages data persistence
   - Handles saving/loading:
     - Caught Pokémon
     - Party members
   - Implements file-based storage with JSON format
   - Provides fallback storage locations

4. **party**
   - Manages Pokémon party system
   - Features:
     - Maximum 6 Pokémon limit
     - Individual Pokémon stats
     - Experience and leveling system
     - Party member management (add/remove)

5. **logger**
   - Provides structured logging
   - Multiple log levels:
     - DEBUG
     - INFO
     - ERROR
     - FATAL
   - Supports multiple output destinations

## Data Models

### Pokemon
```go
type Pokemon struct {
    ID             int
    Name           string
    BaseExperience int
    Height         int
    Weight         int
    Stats          []Stat
    Types          []Type
    Abilities      []Ability
    Moves          []Move
}
```

### PartyPokemon
```go
type PartyPokemon struct {
    InstanceID    string
    Nickname      string
    Level         int
    Experience    int
    CaughtAt      time.Time
    CurrentStats  Stats
    BasePokemon   Pokemon
}
```

## Design Decisions

1. **Caching Strategy**
   - In-memory cache implementation
   - Automatic cleanup to prevent memory leaks
   - Cache invalidation based on time
   - Thread-safe operations for concurrent access

2. **Data Persistence**
   - JSON file-based storage
   - Hierarchical storage location fallback:
     1. Working directory
     2. User's home directory
   - Atomic save operations with mutex locks

3. **Error Handling**
   - Comprehensive error logging
   - Graceful degradation
   - User-friendly error messages
   - Recovery mechanisms for data corruption

4. **Concurrency**
   - Mutex locks for thread safety
   - Goroutines for background tasks
   - Channel-based communication for signals

## Future Extensions

As documented in project_extension_ideas.md, potential improvements include:

1. Command History
   - Up arrow support for previous commands
   - Command history persistence

2. Enhanced Gameplay
   - Pokémon battles
   - Evolution system
   - Random encounters
   - Different types of Pokéballs

3. Technical Improvements
   - Additional unit tests
   - Code refactoring
   - Enhanced error handling
   - Performance optimizations

## Testing

The application includes unit tests, particularly for critical components:
- Cache operations
- API client functionality
- Data persistence
- Party management

## Logging Strategy

Comprehensive logging is implemented across components:

1. **Cache Operations**
   - Cache hits/misses
   - Item addition/removal
   - Cleanup operations

2. **API Operations**
   - Request tracking
   - Rate limiting events
   - Error logging

3. **Pokemon Operations**
   - Catch attempts
   - Party management
   - State changes

4. **Data Persistence**
   - Save/load operations
   - File operations
   - Error handling 