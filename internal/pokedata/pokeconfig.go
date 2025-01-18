package pokedata

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
	"github.com/sakuffo/pokedexcli/internal/pokecache"
)

type Config struct {
	PokeapiClient    pokeapi.Client
	NextLocationsURL *string
	PrevLocationsURL *string
	CaughtPokemon    map[string]pokeapi.Pokemon
	Persistence      *Persistence
	Logger           *logger.Logger
	Party            *party.Party
}

func New(logLevel logger.LogLevel) *Config {

	// Show the log level if its not none
	if logLevel != logger.NONE {
		fmt.Printf("Log level: %v\n", logLevel)
	}

	// Open a log file
	logFile, err := os.OpenFile("pokedexcli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file: %v", err)
	}

	logger := logger.New(logLevel)
	logger.SetWriter(io.MultiWriter(os.Stdout, logFile))

	cache := pokecache.NewCache(5*time.Minute, logger)
	pokeClient := pokeapi.NewClient(5*time.Second, cache, logger)

	persistence, err := NewPersistence("pokedata.json")
	if err != nil {
		logger.Fatal("Failed to initialize persistence: %v", err)
	}
	persistence.SetLogger(logger)

	data, err := persistence.Load()
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Create a new party
	party := &party.Party{
		Members: make([]*party.PartyPokemon, 0),
	}

	// Populate party from loaded data
	if data.PartyMembers != nil {
		party.Members = data.PartyMembers
	}

	cfg := &Config{
		PokeapiClient: pokeClient,
		Persistence:   persistence,
		CaughtPokemon: data.CaughtPokemon,
		Logger:        logger,
		Party:         party,
	}

	return cfg
}
