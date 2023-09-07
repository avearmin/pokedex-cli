package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/avearmin/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	next     string
	previous string
}

func initConfig() config {
	return config{next: pokeapi.LocationAreaEndpoint}
}

func commands() map[string]cliCommand {
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
		"map": {
			name:        "map",
			description: "Gets the next 20 locations.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandHelp(*config) error {
	fmt.Print(
		"Welcome to the Pokedex!\n",
		"Usage:\n\n",
	)
	cmds := commands()
	for key, value := range cmds {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandExit(c *config) error {
	os.Exit(0)
	return nil
}

func commandMap(c *config) error {
	if c.next == "" {
		return fmt.Errorf("No next found")
	}
	data, err := pokeapi.Get(c.next)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func commandMapb(c *config) error {
	if c.previous == "" {
		return fmt.Errorf("No previous found")
	}
	data, err := pokeapi.Get(c.previous)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func printLocationResults(locationResults []pokeapi.LocationResult) {
	for _, result := range locationResults {
		fmt.Println(result.Name)
	}
}

func trimAndLower(s string) string {
	newString := strings.ToLower(s)
	newString = strings.Trim(newString, " ")
	return newString
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := commands()
	config := initConfig()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := trimAndLower(input)
		cmd, ok := cmds[cleanInput]
		if !ok {
			fmt.Printf("%s is not a valid command", cleanInput)
			continue
		}
		err := cmd.callback(&config)
		if err != nil {
			fmt.Println(err)
		}
	}
}
