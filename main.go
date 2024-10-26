package main

import (
	"time"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
	"github.com/sakuffo/pokedexcli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(5 * time.Minute)
	pokeClient := pokeapi.NewClient(5*time.Second, cache)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	startRepl(cfg)
}
