package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	// "github.com/sakuffo/pokedexcli/internal/pokedata" // <-- Remove old import
	"github.com/sakuffo/pokedexcli/internal/commands" // Assuming commands package exists/is correct
	"github.com/sakuffo/pokedexcli/internal/config"   // <-- Add new import
)

// StartRepl starts the Read-Eval-Print Loop for the Pokedex CLI.
func StartRepl(cfg *config.Config) { // <-- Updated type
	scanner := bufio.NewScanner(os.Stdin)
	cmds := commands.GetCommands() // Fetch commands (assuming this function exists and is correct)

	cfg.Logger.Info("Starting REPL...") // Use logger from config

	for {
		fmt.Print("Pokedex > ")
		scanned := scanner.Scan()
		if !scanned {
			// Handle EOF or scanner error
			if err := scanner.Err(); err != nil {
				cfg.Logger.Error("Scanner error: %v", err)
			} else {
				cfg.Logger.Info("EOF received, exiting REPL.")
			}
			break // Exit loop on EOF or error
		}

		text := scanner.Text()
		cleaned := cleanInput(text)
		if len(cleaned) == 0 {
			continue // Skip empty input
		}

		commandName := cleaned[0]
		args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}

		command, ok := cmds[commandName]
		if !ok {
			cfg.Logger.Error("Unknown command: %s", commandName)
			fmt.Println("Unknown command. Type 'help' for available commands.")
			continue
		}

		// Execute the command
		// Command callbacks should ideally return an error for central handling
		err := command.Callback(cfg, args...) // Pass config.Config
		if err != nil {
			// Log the error centrally and inform the user
			cfg.Logger.Error("Command '%s' failed: %v", commandName, err)
			fmt.Printf("Error executing command: %v\n", err)

			// Specific handling for exit command error? Unlikely needed if exit logic is simple.
			if commandName == "exit" {
				// If exit itself fails (e.g., during save triggered by exit),
				// we might already be logging it. Decide if loop should break anyway.
				fmt.Println("Exit command encountered an error, but proceeding to exit.")
				break // Ensure exit happens even if save fails
			}
		}

		// Check if the command was 'exit' to break the loop
		// (Error handling above might also break)
		if commandName == "exit" && err == nil { // Only break if exit command succeeded
			cfg.Logger.Info("Exit command received, terminating REPL.")
			break
		}
	}

	cfg.Logger.Info("REPL loop finished.")
	// Saving is handled in main.go after the REPL exits
}

// cleanInput splits the input string by spaces and trims each part.
func cleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered) // Splits by whitespace and removes empty strings
	return words
}

// Note: Assumes a package `internal/commands` exists with:
// type cliCommand struct { Name string; Description string; Callback func(cfg *config.Config, args ...string) error }
// func GetCommands() map[string]cliCommand
