package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PaginationState struct {
	CurrentURL string
	NextURL    string
	PrevURL    string
}

var pagination = PaginationState{
	CurrentURL: "https://pokeapi.co/api/v2/location/",
	NextURL:    "https://pokeapi.co/api/v2/location/",
	PrevURL:    "",
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type PokeAPIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
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
			callback: func() error {
				response, err := fetchMapData(pagination.CurrentURL)
				if err != nil {
					return err
				}
				fmt.Println(response)
				return nil
			},
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
	if pagination.NextURL == "" {
		fmt.Println("No more pages available")
		return nil
	}
	pagination.CurrentURL = pagination.NextURL
	return displayLocations(pagination.CurrentURL)
}

func commandMapb() error {
	if pagination.PrevURL == "" {
		fmt.Println("No previous page available")
		return nil
	}
	pagination.CurrentURL = pagination.PrevURL
	return displayLocations(pagination.CurrentURL)
}

func displayLocations(url string) error {
	if url == "" {
		fmt.Println("No more pages available")
		return nil
	}

	response, err := fetchMapData(url)
	if err != nil {
		return err
	}

	fmt.Printf("\nDisplaying locations from %s\n", url)
	fmt.Println("Locations:")
	for _, location := range response.Results {
		fmt.Printf("- %s\n", location.Name)
	}

	fmt.Println()

	pagination.NextURL = response.Next
	pagination.PrevURL = response.Previous

	return nil
}

func fetchMapData(url string) (PokeAPIResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching map data: ", err)
		return PokeAPIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status: %s\n", resp.Status)
		return PokeAPIResponse{}, fmt.Errorf("HTTP request failed with status %s", resp.Status)
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
