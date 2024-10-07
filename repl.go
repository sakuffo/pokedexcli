package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Display Next 20 locations in the World",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display Previous 20 locations in the World",
			callback:    commandMapb,
		},
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
	}
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex...")
	os.Exit(0)
	return nil
}

func commandMap() error {
	return nil
}

func commandMapb() error {
	return nil
}

func startRepl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("pokedex > ")

		userInput, err := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if err != nil {
			fmt.Println("Error reading input: ", err)
			continue
		}

		commands := getCommands()

		command, ok := commands[userInput]

		if !ok {
			fmt.Println("Invalid command. Type 'help' to see the list of commands.")
			continue
		}

		err = command.callback()
		if err != nil {
			fmt.Println("Error executing command: ", err)
		}

	}
}
