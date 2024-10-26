package main

import "fmt"

func commandExplore(cfg *config, args []string) error {
	exploreResp, err := cfg.pokeapiClient.FetchAreaPokemon(args[0])
	if err != nil {
		return err
	}

	for _, poke := range exploreResp.PokemonEncounters {
		fmt.Printf("\t- %s\n", poke.Pokemon.Name)
	}

	return nil
}
