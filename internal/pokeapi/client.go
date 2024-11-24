package pokeapi

import (
	"net/http"
	"time"

	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/pokecache"
)

// Client
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache // Keep as pointer
	logger     *logger.Logger
}

func NewClient(timeout time.Duration, cache *pokecache.Cache, logger *logger.Logger) Client {
	logger.Debug("Creating new PokeAPI client")
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:  cache, // Don't dereference
		logger: logger,
	}
}
