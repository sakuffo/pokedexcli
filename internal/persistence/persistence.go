package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sakuffo/pokedexcli/internal/discovery" // Need discovery for Data struct initialization
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// Persistence handles saving and loading application data to/from a file.
type Persistence struct {
	filePath string
	mu       sync.Mutex
	logger   *logger.Logger
}

// NewPersistence creates a new Persistence instance, determining the save file path.
// It prioritizes the working directory and falls back to the user's home directory.
func NewPersistence(filename string) (*Persistence, error) {
	log := logger.New(logger.NONE) // Temporary logger for initialization

	// First try the current working directory
	workDir, err := os.Getwd()
	if err == nil {
		log.Debug("Attempting to use working directory: %s", workDir)

		// Test write permissions
		testFile := filepath.Join(workDir, ".write_test")
		if err := os.WriteFile(testFile, []byte("test"), 0666); err == nil {
			os.Remove(testFile) // Clean up test file
			log.Debug("Write test successful")

			// We have write permission, create directory here
			dir := filepath.Join(workDir, ".pokedexclidata")
			if err := os.MkdirAll(dir, os.ModePerm); err == nil {
				filePath := filepath.Join(dir, filename)
				log.Info("Using save file: %s", filePath)
				// Return with a nil logger initially, expecting SetLogger to be called
				return &Persistence{filePath: filePath, logger: logger.New(logger.NONE)}, nil
			} else {
				log.Error("Failed to create directory: %v", err)
			}
		} else {
			log.Debug("Write test failed: %v", err)
		}
	} else {
		log.Debug("Failed to get working directory: %v", err)
	}

	// Fallback to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	log.Debug("Falling back to home directory: %s", homeDir)
	dir := filepath.Join(homeDir, ".pokedexcli")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory in home: %w", err)
	}

	filePath := filepath.Join(dir, filename)
	log.Info("Using save file: %s", filePath)

	return &Persistence{
		filePath: filePath,
		logger:   logger.New(logger.NONE), // Return with a nil logger initially
	}, nil
}

// Load reads the data from the JSON file.
// If the file doesn't exist, it returns a new empty Data struct.
func (p *Persistence) Load() (*Data, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Info("Save file not found, starting fresh.")
			// Return a new Data struct with initialized maps/slices
			return &Data{
				CaughtPokemon: make(map[string]pokeapi.Pokemon),
				PartyMembers:  make([]*party.PartyPokemon, 0),
				Discoveries:   discovery.NewDiscoveryTracker(), // Initialize runtime tracker
			}, nil
		}
		p.logger.Error("Failed to open save file '%s': %v", p.filePath, err)
		return nil, fmt.Errorf("failed to open save file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields() // Good practice

	var data Data
	err = decoder.Decode(&data)
	if err != nil {
		p.logger.Error("Failed to decode save file '%s': %v", p.filePath, err)
		return nil, fmt.Errorf("failed to decode save file: %w", err)
	}

	p.logger.Debug("Save file loaded successfully from '%s'", p.filePath)

	// Ensure maps/slices are initialized if they were null in the JSON
	if data.CaughtPokemon == nil {
		data.CaughtPokemon = make(map[string]pokeapi.Pokemon)
	}
	if data.PartyMembers == nil {
		data.PartyMembers = make([]*party.PartyPokemon, 0)
	}

	if data.Discoveries == nil {
		data.Discoveries = discovery.NewDiscoveryTracker()
	}

	return &data, nil
}

// Save writes the current data to the JSON file.
// It prepares the data by converting the DiscoveryTracker before saving.
func (p *Persistence) Save(data *Data) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if data.Discoveries == nil {
		data.Discoveries = discovery.NewDiscoveryTracker()
	}

	p.logger.Debug("Attempting to save data to: %s", p.filePath)
	file, err := os.Create(p.filePath)
	if err != nil {
		p.logger.Error("Failed to create/truncate save file '%s': %v", p.filePath, err)
		return fmt.Errorf("failed to create save file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	err = encoder.Encode(data)
	if err != nil {
		p.logger.Error("Failed to encode data to JSON for file '%s': %v", p.filePath, err)
		// Attempt to remove potentially corrupted file? Or leave it? Leaving it for now.
		return fmt.Errorf("failed to encode data: %w", err)
	}

	p.logger.Info("Data saved successfully to %s", p.filePath)
	return nil
}

// SetLogger sets the logger instance for the Persistence object.
func (p *Persistence) SetLogger(lgr *logger.Logger) {
	if lgr != nil {
		p.logger = lgr
	} else {
		// Import the package directly for the constants/functions
		p.logger = logger.New(logger.NONE)
	}
}

// SetLogLevel adjusts the log level of the logger used by Persistence.
func (p *Persistence) SetLogLevel(level logger.LogLevel) {
	p.logger.SetLevel(level)
}
