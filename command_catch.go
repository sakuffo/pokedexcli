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
	randNum := rand.Intn(256)

	if randNum > speciesResp.CaptureRate {
		fmt.Printf("%s escaped!\n", speciesResp.Name)
	} else {
		fmt.Printf("%s was caught!\n", speciesResp.Name)
		err := cfg.AddToPokedex(pokemonName)
		if err != nil {
			return err
		}
	}

	return nil
}
