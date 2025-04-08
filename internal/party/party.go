package party

import (
	"errors"
	"fmt"
	"sync"
)

type Party struct {
	Members []*PartyPokemon
	mu      sync.Mutex
}

func (p *Party) AddMember(pokemon *PartyPokemon) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Members) >= 6 {
		return errors.New("party is full (maximum 6 Pokemon)")
	}

	for _, member := range p.Members {
		if member.BasePokemon.Name == pokemon.BasePokemon.Name {
			return fmt.Errorf("Pokemon %s is already in the party", pokemon.BasePokemon.Name)
		}
	}

	p.Members = append(p.Members, pokemon)
	return nil
}

func (p *Party) RemoveMember(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, member := range p.Members {
		if member.BasePokemon.Name == name {
			p.Members[i] = p.Members[len(p.Members)-1]
			p.Members = p.Members[:len(p.Members)-1]
			return nil
		}
	}
	return fmt.Errorf("Pokemon %s not found in party", name)
}

func (p *Party) ListMembers() []*PartyPokemon {
	p.mu.Lock()
	defer p.mu.Unlock()

	result := make([]*PartyPokemon, len(p.Members))
	copy(result, p.Members)
	return result
}

func (p *Party) GetMember(name string) (*PartyPokemon, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, member := range p.Members {
		if member.BasePokemon.Name == name {
			return member, true
		}
	}
	return nil, false
}

func (p *Party) IsFull() bool {
	return len(p.Members) >= 6
}
