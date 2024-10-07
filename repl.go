package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
