package app

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sakuffo/pokedexcli/internal/cache"
	"github.com/sakuffo/pokedexcli/internal/config"
	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/persistence"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// Initialize sets up all application components and returns the config.
// It returns any error encountered rather than fatal-exiting.
func Initialize(logLevel logger.LogLevel) (*config.Config, error) {
	// Set up logging
	logFile, err := setupLogFile("pokedexcli.log")
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	appLogger := logger.New(logLevel)
	configureLogger(appLogger, logFile, logLevel)

	// Initialize components
	appCache := cache.NewCache(5*time.Minute, appLogger)
	pokeClient := pokeapi.NewClient(5*time.Second, appCache, appLogger)

	persister, err := persistence.NewPersistence("pokedata.json")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize persistence: %w", err)
	}
	persister.SetLogger(appLogger)

	// Load data
	loadedData, err := persister.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// Ensure proper initialization of components
	discoveryTracker := ensureDiscoveryTracker(loadedData.Discoveries)
	partyManager := setupParty(loadedData.PartyMembers)

	// Create application state
	cfg := &config.Config{
		PokeapiClient: pokeClient,
		Persistence:   persister,
		CaughtPokemon: loadedData.CaughtPokemon,
		Discoveries:   discoveryTracker,
		Logger:        appLogger,
		Party:         partyManager,
	}

	appLogger.Info("Application initialized successfully")
	return cfg, nil
}

// SaveData saves the current application state.
func SaveData(cfg *config.Config) error {
	if cfg.Persistence == nil {
		return fmt.Errorf("persistence service not initialized")
	}

	dataToSave := &persistence.Data{
		CaughtPokemon: cfg.CaughtPokemon,
		PartyMembers:  cfg.Party.Members,
		Discoveries:   cfg.Discoveries,
	}

	return cfg.Persistence.Save(dataToSave)
}

// Helper functions
func setupLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func configureLogger(lgr *logger.Logger, logFile *os.File, level logger.LogLevel) {
	if level == logger.NONE {
		lgr.SetWriter(logFile)
	} else {
		lgr.SetWriter(io.MultiWriter(os.Stdout, logFile))
	}
}

func ensureDiscoveryTracker(tracker *discovery.DiscoveryTracker) *discovery.DiscoveryTracker {
	if tracker != nil && tracker.IsInitialized() {
		return tracker
	}
	return discovery.NewDiscoveryTracker()
}

func setupParty(members []*party.PartyPokemon) *party.Party {
	p := &party.Party{
		Members: make([]*party.PartyPokemon, 0),
	}
	if members != nil {
		p.Members = members
	}
	return p
}
