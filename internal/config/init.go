package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sakuffo/pokedexcli/internal/cache"
	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/persistence"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// New creates and initializes the application configuration and state.
// This is a temporary approach until initialization is moved to main.go.
func New(logLevel logger.LogLevel) *Config {
	// Show the log level if it's not none
	if logLevel != logger.NONE {
		fmt.Printf("Log level: %v\n", logLevel)
	}

	// Open a log file
	logFile, err := os.OpenFile("pokedexcli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	appLogger := logger.New(logLevel)
	// Configure logger output
	if logLevel == logger.NONE {
		// If NONE, only log to file, not stdout
		appLogger.SetWriter(logFile)
	} else {
		appLogger.SetWriter(io.MultiWriter(os.Stdout, logFile))
	}

	appCache := cache.NewCache(5*time.Minute, appLogger)
	pokeClient := pokeapi.NewClient(5*time.Second, appCache, appLogger)

	// Use the new persistence package
	persister, err := persistence.NewPersistence("pokedata.json")
	if err != nil {
		appLogger.Fatal("Failed to initialize persistence: %v", err)
	}
	persister.SetLogger(appLogger)

	// Load data using the persister
	loadedData, err := persister.Load()
	if err != nil {
		appLogger.Fatal("Failed to load data: %v", err)
	}

	// Ensure discovery tracker is properly initialized
	var tracker *discovery.DiscoveryTracker
	if loadedData.Discoveries != nil && loadedData.Discoveries.IsInitialized() {
		tracker = loadedData.Discoveries
	} else {
		// If not properly initialized, create a new one
		tracker = discovery.NewDiscoveryTracker()
	}

	// Create a new party
	currentParty := &party.Party{
		Members: make([]*party.PartyPokemon, 0),
	}

	// Populate party from loaded data
	if loadedData.PartyMembers != nil {
		currentParty.Members = loadedData.PartyMembers
		appLogger.Debug("Loaded %d members into party", len(currentParty.Members))
	}

	// Create the config struct
	cfg := &Config{
		PokeapiClient: pokeClient,
		Persistence:   persister,
		CaughtPokemon: loadedData.CaughtPokemon,
		Discoveries:   tracker, // Use the properly initialized tracker
		Logger:        appLogger,
		Party:         currentParty,
	}

	appLogger.Info("Application initialized successfully.")
	return cfg
}
