package party

import (
	"errors"
	"sync"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

type Party struct {
	Members []pokeapi.Pokemon
	mu      sync.Mutex
}

func (p *Party) AddMember(pokemon pokeapi.Pokemon) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Members) >= 6 {
		return errors.New("party is full")
	}

	for _, member := range p.Members {
		if member.ID == pokemon.ID {
			return errors.New("pokemon is already in the party")
		}
	}

	p.Members = append(p.Members, pokemon)
	return nil
}

func (p *Party) RemoveMember(pokemonName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, member := range p.Members {
		if member.Name == pokemonName {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			return nil
		}
	}
	return errors.New("pokemon not found in party")
}

func (p *Party) ListMembers() []pokeapi.Pokemon {
	p.mu.Lock()
	defer p.mu.Unlock()

	return append([]pokeapi.Pokemon(nil), p.Members...)
}
