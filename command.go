package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type PokeAPIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
		"test": {
			name:        "test",
			description: "Test the Pokedex",
			callback:    fetchMapData,
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

func fetchMapData() (PokeAPIResponse, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/location/")
	if err != nil {
		fmt.Println("Error fetching map data: ", err)
		return PokeAPIResponse{}, err
	}
	var mapData PokeAPIResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading map data: ", err)
		return PokeAPIResponse{}, err
	}

	err = json.Unmarshal(body, &mapData)
	if err != nil {
		fmt.Println("Error unmarshalling map data: ", err)
		return PokeAPIResponse{}, err
	}

	return mapData, nil
}
