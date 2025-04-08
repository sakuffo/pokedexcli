package config

import (
	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/persistence"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// Config holds the runtime state of the application but doesn't initialize it.
type Config struct {
	PokeapiClient    pokeapi.Client
	NextLocationsURL *string
	PrevLocationsURL *string
	CaughtPokemon    map[string]pokeapi.Pokemon
	Discoveries      *discovery.DiscoveryTracker
	Persistence      *persistence.Persistence
	Logger           *logger.Logger
	Party            *party.Party
}
