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
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Lists the next 20 locations",
			Callback:    commandMapf,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Lists the previous 20 locations",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Lists the pokemon in the area",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempts to catch a pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspects a pokemon in you have caught",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists all the pokemon you have caught",
			Callback:    commandPokedex,
		},
		"party": {
			Name:        "party",
			Description: "Lists all the pokemon in your party",
			Callback:    commandParty,
		},
	}
}
