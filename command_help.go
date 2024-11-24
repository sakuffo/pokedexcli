package main

import "fmt"

func commandHelp(cfg *config, args ...string) error {
	cfg.logger.Debug("Executing 'help' command")
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	cfg.logger.Debug("Displayed help")
	return nil
}
