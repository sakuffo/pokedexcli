package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// types

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.PokeAPIPokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

// functions

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lists the next 20 locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists the pokemon in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a pokemon in your pokedex",
			callback:    commandInspect,
		},
	}
}

func (c *config) AddToPokedex(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	pokedex := c.caughtPokemon

	pokemonResp, err := c.pokeapiClient.FetchPokemon(name)
	if err != nil {
		fmt.Println("Error fetching pokemon:", err)
		return err
	}

	pokedex[pokemonResp.Name] = pokemonResp

	return nil
}

func startRepl(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Invalid command")
			continue
		}
	}
}
