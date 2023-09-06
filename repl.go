package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandMap() map[string]cliCommand {
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
	cmdMap := commandMap()
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

func trimAndLower(s string) string {
	newString := strings.ToLower(s)
	newString = strings.Trim(newString, " ")
	return newString
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cmdMap := commandMap()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cmd, ok := cmdMap[trimAndLower(input)]
		if !ok {
			fmt.Printf("%s is not a valid command", input)
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}
