package main

import (
	"time"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
	"github.com/sakuffo/pokedexcli/internal/pokecache"
	"golang.org/x/exp/rand"
)

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	pokedex := pokeapi.PokedexPokeInfo{
		Pokemon: make(map[string]pokeapi.PokeAPIPokemon),
	}
	cache := pokecache.NewCache(5 * time.Minute)
	pokeClient := pokeapi.NewClient(5*time.Second, cache)
	cfg := &config{
		pokeapiClient: pokeClient,
		pokedex:       pokedex,
	}

	startRepl(cfg)
}
