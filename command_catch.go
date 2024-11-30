package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/rand"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		cfg.logger.Error("Pokemon name is required")
		return errors.New("pokemon name is required")
	}

	pokemonName := args[0]
	cfg.logger.Debug("Attempting to catch: %s", pokemonName)

	pokemonResp, err := cfg.pokeapiClient.FetchPokemon(pokemonName)
	if err != nil {
		cfg.logger.Error("Failed to fetch pokemon %s: %v", pokemonName, err)
		return err
	}

	res := rand.Intn(pokemonResp.BaseExperience)
	cfg.logger.Debug("Catch roll: %d vs threshold: 40 (base exp: %d)", res, pokemonResp.BaseExperience) // Add log

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if res > 40 {
		cfg.logger.Debug("%s escaped!", pokemonName)
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	cfg.logger.Debug("Successfully caught %s!", pokemonName)
	fmt.Printf("%s was caught!\n", pokemonResp.Name)
	fmt.Println("You may now inspect it using the inspect command")

	cfg.caughtPokemon[pokemonResp.Name] = pokemonResp

	err = cfg.party.AddMember(pokemonResp)
	if err != nil {
		cfg.logger.Error("Failed to add %s to party: %v", pokemonName, err)
		fmt.Printf("Failed to add %s to party: %v\n", pokemonName, err)
	} else {
		fmt.Printf("%s was added to your party.\n", pokemonName)
	}

	err = saveData(cfg)
	if err != nil {
		cfg.logger.Error("Failed to save data after catching %s: %v", pokemonName, err)
		fmt.Printf("Failed to save data: %v\n", err)
	}

	return nil
}
