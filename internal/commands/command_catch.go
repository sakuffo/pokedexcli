package commands

import (
	"errors"
	"fmt"

	"golang.org/x/exp/rand"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func commandCatch(cfg *pokeconfig.Config, args ...string) error {
	if len(args) != 1 {
		cfg.Logger.Error("Pokemon name is required")
		return errors.New("pokemon name is required")
	}

	pokemonName := args[0]
	cfg.Logger.Debug("Attempting to catch: %s", pokemonName)

	pokemonResp, err := cfg.PokeapiClient.FetchPokemon(pokemonName)
	if err != nil {
		cfg.Logger.Error("Failed to fetch pokemon %s: %v", pokemonName, err)
		return err
	}

	res := rand.Intn(pokemonResp.BaseExperience)
	cfg.Logger.Debug("Catch roll: %d vs threshold: 40 (base exp: %d)", res, pokemonResp.BaseExperience) // Add log

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if res > 40 {
		cfg.Logger.Debug("%s escaped!", pokemonName)
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	cfg.Logger.Debug("Successfully caught %s!", pokemonName)
	fmt.Printf("%s was caught!\n", pokemonResp.Name)
	fmt.Println("You may now inspect it using the inspect command")

	cfg.CaughtPokemon[pokemonResp.Name] = pokemonResp

	err = cfg.Party.AddMember(pokemonResp)
	if err != nil {
		cfg.Logger.Error("Failed to add %s to party: %v", pokemonName, err)
		fmt.Printf("Failed to add %s to party: %v\n", pokemonName, err)
	} else {
		fmt.Printf("%s was added to your party.\n", pokemonName)
	}

	err = pokedata.SaveData(cfg)
	if err != nil {
		cfg.Logger.Error("Failed to save data after catching %s: %v", pokemonName, err)
		fmt.Printf("Failed to save data: %v\n", err)
	}

	return nil
}
