package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
	"github.com/sakuffo/pokedexcli/internal/pokecache"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

// types

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.PokeAPIPokemon
	persistence      *pokedata.Persistence
	logger           *logger.Logger
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
			description: "Inspects a pokemon in you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}

func InitializeConfig(logLevel logger.LogLevel) *config {

	// Show the log level if its not none
	if logLevel != logger.NONE {
		fmt.Printf("Log level: %v\n", logLevel)
	}

	logFile, err := os.OpenFile("pokedexcli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger := logger.New(logLevel)
	logger.SetWriter(io.MultiWriter(os.Stdout, logFile))

	cache := pokecache.NewCache(5 * time.Minute)
	pokeClient := pokeapi.NewClient(5*time.Second, cache)

	persistence, err := pokedata.NewPersistence("pokedata.json")
	if err != nil {
		logger.Fatal("Failed to initialize persistence: %v", err)
	}
	persistence.SetLogger(logger)

	data, err := persistence.Load()
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	cfg := &config{
		pokeapiClient: pokeClient,
		persistence:   persistence,
		caughtPokemon: data.CaughtPokemon,
		logger:        logger,
	}

	return cfg
}

func saveData(cfg *config) error {
	if cfg.persistence == nil {
		return errors.New("persistence not initialized")
	}

	data := &pokedata.Data{
		CaughtPokemon: cfg.caughtPokemon,
	}

	err := cfg.persistence.Save(data)
	if err != nil {
		return err
	}

	return nil
}

func startRepl(cfg *config) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, saving data and exiting...")
		err := saveData(cfg)
		if err != nil {
			fmt.Printf("Failed to save data: %v\n", err)
		}
		os.Exit(0)
	}()

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
