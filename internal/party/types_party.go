package party

import (
	"time"

	"github.com/google/uuid"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

type PartyPokemon struct {
	// Instance-specific fields
	InstanceID string    `json:"instance_id"`
	Nickname   string    `json:"nickname"`
	Level      int       `json:"level"`
	Experience int       `json:"experience"`
	CaughtAt   time.Time `json:"caught_at"`

	// Current stats (calculated from base stats)
	CurrentStats struct {
		HP             int `json:"hp"`
		Attack         int `json:"attack"`
		Defense        int `json:"defense"`
		SpecialAttack  int `json:"special_attack"`
		SpecialDefense int `json:"special_defense"`
		Speed          int `json:"speed"`
	} `json:"current_stats"`

	// Reference to base Pokemon
	BasePokemon pokeapi.Pokemon `json:"base_pokemon"`
}

func NewPartyPokemon(base pokeapi.Pokemon) *PartyPokemon {
	return &PartyPokemon{
		InstanceID:   uuid.New().String(), // Generate a new UUID
		Nickname:     base.Name,
		Level:        5,
		Experience:   0,
		CaughtAt:     time.Now(),
		BasePokemon:  base,
		CurrentStats: calculateInitialStats(base),
	}
}

func calculateInitialStats(base pokeapi.Pokemon) struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
} {
	var stats struct {
		HP             int `json:"hp"`
		Attack         int `json:"attack"`
		Defense        int `json:"defense"`
		SpecialAttack  int `json:"special_attack"`
		SpecialDefense int `json:"special_defense"`
		Speed          int `json:"speed"`
	}

	for _, stat := range base.Stats {
		switch stat.Stat.Name {
		case "hp":
			stats.HP = stat.BaseStat
		case "attack":
			stats.Attack = stat.BaseStat
		case "defense":
			stats.Defense = stat.BaseStat
		case "special-attack":
			stats.SpecialAttack = stat.BaseStat
		case "special-defense":
			stats.SpecialDefense = stat.BaseStat
		case "speed":
			stats.Speed = stat.BaseStat
		}
	}

	return stats
}
