package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sakuffo/pokedexcli/internal/commands"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

// functions

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func startRepl(cfg *pokedata.Config) {
	cfg.Logger.Debug("Starting REPL")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cfg.Logger.Debug("Received interrupt signal, saving data and exiting...")
		fmt.Println("\nReceived interrupt signal, saving data and exiting...")
		err := pokedata.SaveData(cfg)
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

		command, exists := commands.GetCommands()[commandName]
		if exists {
			err := command.Callback(cfg, args...)
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
