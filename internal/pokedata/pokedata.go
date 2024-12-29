package pokedata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

type Data struct {
	CaughtPokemon map[string]pokeapi.Pokemon `json:"caught_pokemon"`
	PartyMembers  []*party.PartyPokemon      `json:"party_members"`
}
type Persistence struct {
	filePath string
	mu       sync.Mutex
	logger   *logger.Logger
}

func NewPersistence(filename string) (*Persistence, error) {
	log := logger.New(logger.NONE)

	// First try the current working directory
	workDir, err := os.Getwd()
	if err == nil {
		log.Debug("Attempting to use working directory: %s\n", workDir)

		// Test write permissions
		testFile := filepath.Join(workDir, ".write_test")
		if err := os.WriteFile(testFile, []byte("test"), 0666); err == nil {
			os.Remove(testFile) // Clean up test file
			log.Debug("Write test successful")

			// We have write permission, create directory here
			dir := filepath.Join(workDir, ".pokedexclidata")
			if err := os.MkdirAll(dir, os.ModePerm); err == nil {
				filePath := filepath.Join(dir, filename)
				log.Info("Using save file: %s\n", filePath)
				return &Persistence{filePath: filePath}, nil
			} else {
				log.Error("Failed to create directory: %v\n", err)
			}
		} else {
			log.Debug("Write test failed: %v\n", err)
		}
	}

	// Fallback to home directory if executable directory isn't writable
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get both executable and home directory: %w", err)
	}

	log.Debug("Falling back to home directory: %s\n", homeDir)
	dir := filepath.Join(homeDir, ".pokedexcli")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory in home: %w", err)
	}

	filePath := filepath.Join(dir, filename)
	log.Info("Using save file: %s\n", filePath)

	return &Persistence{
		filePath: filePath,
		logger:   log,
	}, nil
}

func (p *Persistence) Load() (*Data, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Data{CaughtPokemon: make(map[string]pokeapi.Pokemon)}, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	var data Data
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.CaughtPokemon == nil {
		data.CaughtPokemon = make(map[string]pokeapi.Pokemon)
	}

	if data.PartyMembers == nil {
		data.PartyMembers = []*party.PartyPokemon{}
	}

	return &data, nil
}

func (p *Persistence) SetLogger(logger *logger.Logger) {
	p.logger = logger
}
func (p *Persistence) Save(data *Data) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Debug("Attempting to save to: %s", p.filePath)
	file, err := os.Create(p.filePath)
	if err != nil {
		p.logger.Error("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		p.logger.Error("Failed to encode data: %v", err)
		return err
	}

	p.logger.Info("Save successful")
	return nil
}

func (p *Persistence) SetLogLevel(level logger.LogLevel) {
	p.logger.SetLevel(level)
}
