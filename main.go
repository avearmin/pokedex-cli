package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func CommandMap() map[string]cliCommand {
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
	}
}

func commandHelp() error {
	fmt.Print(
		"Welcome to the Pokedex!\n",
		"Usage:\n\n",
	)
	cmdMap := CommandMap()
	for key, value := range cmdMap {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmdMap := CommandMap()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cmd, ok := cmdMap[userInput]
		if !ok {
			fmt.Printf("%s is not a valid command", userInput)
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}
