package main

import (
	"fmt"
	"os"
	"time"

	"github.com/avearmin/pokedex-cli/internal/inputparser"
	"github.com/avearmin/pokedex-cli/internal/pokeapi"
	"github.com/avearmin/pokedex-cli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, ca *pokecache.Cache, arg string) error
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
		"explore": {
			name:        "explore <area-name>",
			description: "Gets a list of pokemon in the location",
			callback:    commandExplore,
		},
	}
}

func commandHelp(c *config, cache *pokecache.Cache, arg string) error {
	fmt.Print(
		"Welcome to the Pokedex!\n",
		"Usage:\n\n",
	)
	cmds := commands()
	for _, value := range cmds {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Println()
	return nil
}

func commandExit(c *config, cache *pokecache.Cache, arg string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *config, cache *pokecache.Cache, arg string) error {
	if c.next == "" {
		return fmt.Errorf("No next found")
	}
	data, err := pokeapi.Get(c.next, pokeapi.LocationData{}, cache)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func commandMapb(c *config, cache *pokecache.Cache, arg string) error {
	if c.previous == "" {
		return fmt.Errorf("No previous found")
	}
	data, err := pokeapi.Get(c.previous, pokeapi.LocationData{}, cache)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func commandExplore(c *config, cache *pokecache.Cache, arg string) error {
	data, err := pokeapi.Get(pokeapi.LocationAreaEndpoint+arg, pokeapi.LocationAreaData{}, cache)
	if err != nil {
		return err
	}
	printPokemonNames(data.PokemonEncounters)
	return nil

}

func printLocationResults(locationResults []pokeapi.LocationResult) {
	for _, result := range locationResults {
		fmt.Println(result.Name)
	}
}

func printPokemonNames(encounters []pokeapi.PokemonEncounters) {
	for _, result := range encounters {
		fmt.Println(result.Pokemon.Name)
	}
}

func StartRepl() {
	inputParser := inputparser.NewInputParser(2)
	cmds := commands()
	config := initConfig()
	cache := pokecache.NewCache(5 * time.Minute)
	for {
		fmt.Print("Pokedex > ")
		if err := inputParser.ScanAndParse(); err != nil {
			fmt.Println(err)
		}
		cmd, ok := cmds[inputParser.Arg(0)]
		if !ok {
			fmt.Printf("%s is not a valid command\n", inputParser.Arg(0))
			continue
		}
		err := cmd.callback(&config, cache, inputParser.Arg(1))
		if err != nil {
			fmt.Println(err)
		}
	}
}
