package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("pokemon name is required")
	}

	pokemonName := args[0]
	speciesResp, err := cfg.pokeapiClient.FetchPokemonSpecies(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", speciesResp.Name)

	rand.Seed(uint64(time.Now().UnixNano()))
	randNum := rand.Intn(100)

	if randNum > speciesResp.CaptureRate {
		fmt.Printf("You failed to catch the pokemon\n")
	} else {
		fmt.Printf("You caught the pokemon\n")
	}

	return nil
}
