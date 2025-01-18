package commands

import "github.com/sakuffo/pokedexcli/internal/pokedata"

type cliCommand struct {
	Name        string
	Description string
	Callback    func(*pokedata.Config, ...string) error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Lists the next 20 locations",
			Callback:    CommandMapf,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Lists the previous 20 locations",
			Callback:    CommandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Lists the pokemon in the area",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempts to catch a pokemon",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspects a pokemon in you have caught",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists all the pokemon you have caught",
			Callback:    CommandPokedex,
		},
		"party": {
			Name:        "party",
			Description: "Lists all the pokemon in your party",
			Callback:    CommandParty,
		},
	}
}
