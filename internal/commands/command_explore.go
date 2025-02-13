package commands

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func CommandExplore(cfg *pokedata.Config, args ...string) error {

	// Check if area is provided
	if len(args) != 1 {
		cfg.Logger.Error("Area is required")
		return errors.New("area is required")
	}

	area := args[0]

	cfg.Logger.Info("Exploring area: %s", area)

	exploreResp, err := cfg.PokeapiClient.FetchAreaPokemon(area)
	if err != nil {
		cfg.Logger.Error("Failed to fetch area %s: %v", area, err)
		return err
	}
	locationName := exploreResp.Location.Name

	// Initialize the area in discovered map if it doesn't exist
	if _, exists := cfg.DiscoveredPokemon[locationName]; !exists {
		cfg.DiscoveredPokemon[locationName] = make(map[string]bool)
	}

	// Get all undiscovered pokemon in this area
	var undiscovered []string
	discoveredInArea := 0
	totalInArea := len(exploreResp.PokemonEncounters)

	for _, poke := range exploreResp.PokemonEncounters {
		if !cfg.DiscoveredPokemon[locationName][poke.Pokemon.Name] {
			undiscovered = append(undiscovered, poke.Pokemon.Name)
		} else {
			discoveredInArea++
		}
	}

	fmt.Printf("Exploring %s...\n\n", area)

	// Track newly discovered pokemon for coloring
	newlyDiscovered := make(map[string]bool)

	if len(undiscovered) > 0 {
		// Randomly discover 1-3 new pokemon
		numToDiscover := rand.Intn(3) + 1
		if numToDiscover > len(undiscovered) {
			numToDiscover = len(undiscovered)
		}

		// Shuffle the undiscovered slice
		rand.Shuffle(len(undiscovered), func(i, j int) {
			undiscovered[i], undiscovered[j] = undiscovered[j], undiscovered[i]
		})

		// Mark the newly discovered pokemon
		fmt.Println("\nYou discovered new Pokemon!\n")
		for i := 0; i < numToDiscover; i++ {
			pokeName := undiscovered[i]
			cfg.DiscoveredPokemon[locationName][pokeName] = true
			newlyDiscovered[pokeName] = true
			discoveredInArea++
		}
	}

	// Show all Pokemon in this area and their discovery status
	fmt.Printf("Pokemon in %s:\n", area)
	for _, poke := range exploreResp.PokemonEncounters {
		if cfg.DiscoveredPokemon[locationName][poke.Pokemon.Name] {
			if newlyDiscovered[poke.Pokemon.Name] {
				fmt.Printf("\t- %s%s%s (New!)\n", colorGreen, poke.Pokemon.Name, colorReset)
			} else {
				fmt.Printf("\t- %s \n", poke.Pokemon.Name)
			}
		} else {
			fmt.Printf("\t- ??? \n")
		}
	}

	fmt.Printf("\nProgress for this area: %d/%d Pokemon discovered\n", discoveredInArea, totalInArea)
	return nil
}
