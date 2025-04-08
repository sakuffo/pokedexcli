package commands

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/sakuffo/pokedexcli/internal/config"
)

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func CommandExplore(cfg *config.Config, args ...string) error {

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

	// Discovery Logic
	var undiscovered []string
	totalInArea := len(exploreResp.PokemonEncounters)

	for _, poke := range exploreResp.PokemonEncounters {
		if !cfg.Discoveries.IsDiscovered(locationName, poke.Pokemon.Name) {
			undiscovered = append(undiscovered, poke.Pokemon.Name)
		}
	}

	discoveredInArea := cfg.Discoveries.CountDiscoveredInLocation(locationName)

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
		fmt.Println("\nYou discovered new Pokemon!")
		for i := 0; i < numToDiscover; i++ {
			pokeName := undiscovered[i]
			cfg.Logger.Info("Marking %s as discovered in %s", pokeName, locationName)
			cfg.Discoveries.MarkDiscovered(locationName, pokeName)

			cfg.Logger.Info("Discovered %s in %s", pokeName, locationName)
			fmt.Printf("Discovered %s in %s\n", pokeName, locationName)

			newlyDiscovered[pokeName] = true
			discoveredInArea++
		}
	}

	cfg.Logger.Info("Progress for this area: %d/%d Pokemon discovered", discoveredInArea, totalInArea)
	// List all discovered pokemon in this area
	discovered := cfg.Discoveries.GetDiscoveredInLocation(locationName)
	for _, poke := range discovered {
		// Print the new pokemon in green and the rest in reset
		if newlyDiscovered[poke] {
			fmt.Printf("\t- %s%s%s\n", colorGreen, poke, colorReset)
		} else {
			fmt.Printf("\t- %s\n", poke)
		}
	}

	fmt.Printf("\nProgress for this area: %d/%d Pokemon discovered\n", discoveredInArea, totalInArea)
	return nil
}
