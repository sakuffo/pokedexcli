package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

// types

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeconfig.Config, ...string) error
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
			description: "Inspects a pokemon in you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon you have caught",
			callback:    commandPokedex,
		},
		"party": {
			name:        "party",
			description: "Lists all the pokemon in your party",
			callback:    commandParty,
		},
	}
}

func saveData(cfg *pokeconfig.Config) error {
	cfg.Logger.Debug("Saving data")
	if cfg.Persistence == nil {
		cfg.Logger.Error("Persistence not initialized")
		return errors.New("persistence not initialized")
	}

	data := &pokedata.Data{
		CaughtPokemon: cfg.CaughtPokemon,
		PartyMembers:  cfg.Party.Members,
	}

	err := cfg.Persistence.Save(data)
	if err != nil {
		cfg.Logger.Error("Failed to save data: %v", err)
		return err
	}

	return nil
}

func startRepl(cfg *pokeconfig.Config) {
	cfg.Logger.Debug("Starting REPL")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cfg.Logger.Debug("Received interrupt signal, saving data and exiting...")
		fmt.Println("\nReceived interrupt signal, saving data and exiting...")
		err := saveData(cfg)
		if err != nil {
			cfg.Logger.Error("Failed to save data: %v", err)
			fmt.Printf("Failed to save data: %v\n", err)
		}
		os.Exit(0)
	}()

	reader := bufio.NewScanner(os.Stdin)
	for {
		cfg.Logger.Debug("Getting input")
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

		cfg.Logger.Debug("User entered command: %s with args: %v", commandName, args)

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				cfg.Logger.Error("Error executing command: '%s': %v", commandName, err)
				fmt.Println(err)
			}
		} else {
			cfg.Logger.Error("Invalid command: '%s'", commandName)
			fmt.Println("Invalid command")
			continue
		}
	}
}
