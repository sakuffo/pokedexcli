package pokeapi

import (
	"net/http"
	"time"

	"github.com/sakuffo/pokedexcli/internal/cache"
	"github.com/sakuffo/pokedexcli/internal/logger"
)

// Client
type Client struct {
	httpClient http.Client
	cache      *cache.Cache // Keep as pointer
	logger     *logger.Logger
}

func NewClient(timeout time.Duration, cache *cache.Cache, logger *logger.Logger) Client {
	logger.Debug("Creating new PokeAPI client")
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:  cache, // Don't dereference
		logger: logger,
	}
}
