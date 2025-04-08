package persistence

import (
	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// Data defines the structure for saving and loading application state.
type Data struct {
	CaughtPokemon map[string]pokeapi.Pokemon  `json:"caught_pokemon"`
	PartyMembers  []*party.PartyPokemon       `json:"party_members"`
	Discoveries   *discovery.DiscoveryTracker `json:"discoveries"`
}
