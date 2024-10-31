package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/rand"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("pokemon name is required")
	}

	pokemonName := args[0]
	pokemonResp, err := cfg.pokeapiClient.FetchPokemon(pokemonName)
	if err != nil {
		fmt.Printf("Error fetching pokemon: %s\n", err)
		return err
	}

	res := rand.Intn(pokemonResp.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonResp.Name)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemonResp.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonResp.Name)
	fmt.Println("You may now inspect it using the inspect command")

	cfg.caughtPokemon[pokemonResp.Name] = pokemonResp
	return nil
}
