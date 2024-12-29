package party

import (
	"errors"
	"sync"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

type Party struct {
	Members []*PartyPokemon
	mu      sync.Mutex
}

func (p *Party) AddMember(pokemon pokeapi.Pokemon) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Members) >= 6 {
		return errors.New("party is full")
	}

	partyPokemon := NewPartyPokemon(pokemon)

	for _, member := range p.Members {
		if member.BasePokemon.ID == pokemon.ID {
			return errors.New("pokemon is already in the party") // TODO: Why can't we use the same pokemon in the party?
		}
	}

	p.Members = append(p.Members, partyPokemon)
	return nil
}

func (p *Party) RemoveMember(nickname string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, member := range p.Members {
		if member.Nickname == nickname {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			return nil
		}
	}
	return errors.New("pokemon not found in party")
}

func (p *Party) ListMembers() []*PartyPokemon {
	p.mu.Lock()
	defer p.mu.Unlock()

	result := make([]*PartyPokemon, len(p.Members))
	copy(result, p.Members)
	return result
}

func (p *Party) GetMember(nickname string) (*PartyPokemon, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, member := range p.Members {
		if member.Nickname == nickname {
			return member, nil
		}
	}
	return nil, errors.New("pokemon not found in party")
}
