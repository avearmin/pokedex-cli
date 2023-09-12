package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/avearmin/pokedex-cli/internal/pokeapi"
	"github.com/avearmin/pokedex-cli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args []string) error
}

type config struct {
	cache    *pokecache.Cache
	pokedex  map[string]pokeapi.PokemonData
	next     string
	previous string
}

func initConfig() config {
	cache := pokecache.NewCache(5 * time.Minute)
	pokedex := make(map[string]pokeapi.PokemonData)
	return config{
		cache:   cache,
		pokedex: pokedex,
		next:    pokeapi.LocationAreaEndpoint,
	}
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
		"catch": {
			name:        "catch <pokemon>",
			description: "Attempt to catch a pokemon and save them to your pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Inspect a pokemon in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show a list of all pokemon in your pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(c *config, args []string) error {
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

func commandExit(c *config, args []string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *config, args []string) error {
	if c.next == "" {
		return fmt.Errorf("No next found")
	}
	data, err := pokeapi.Get(c.next, pokeapi.LocationData{}, c.cache)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func commandMapb(c *config, args []string) error {
	if c.previous == "" {
		return fmt.Errorf("No previous found")
	}
	data, err := pokeapi.Get(c.previous, pokeapi.LocationData{}, c.cache)
	if err != nil {
		return err
	}
	c.previous = data.Previous
	c.next = data.Next
	printLocationResults(data.Results)
	return nil
}

func commandExplore(c *config, args []string) error {
	data, err := pokeapi.Get(pokeapi.LocationAreaEndpoint+args[0], pokeapi.LocationAreaData{}, c.cache)
	if err != nil {
		return err
	}
	printPokemonNames(data.PokemonEncounters)
	return nil
}

func commandCatch(c *config, args []string) error {
	data, err := pokeapi.Get(pokeapi.PokemonEndpoint+args[0], pokeapi.PokemonData{}, c.cache)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing pokeball at %s...\n", args[0])
	if chanceToCatch := rand.Intn(data.BaseExperience); chanceToCatch > 50 {
		fmt.Printf("%s got away\n", args[0])
		return nil
	}
	fmt.Printf("Caught %s!\n", args[0])
	c.pokedex[args[0]] = data
	return nil
}

func commandInspect(c *config, args []string) error {
	data, ok := c.pokedex[args[0]]
	if !ok {
		return fmt.Errorf("%s was not found in your Pokedex\n", args[0])
	}
	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range data.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range data.Types {
		fmt.Printf("-%s\n", pokemonType.Type.Name)
	}
	return nil
}

func commandPokedex(c *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for pokemon, _ := range c.pokedex {
		fmt.Printf("-%s\n", pokemon)
	}
	return nil
}

func trimAndLower(s string) string {
	newString := strings.ToLower(s)
	newString = strings.Trim(newString, " ")
	return newString
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
	scanner := bufio.NewScanner(os.Stdin)
	cmds := commands()
	config := initConfig()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := trimAndLower(scanner.Text())
		args := strings.Fields(input)
		cmd, ok := cmds[args[0]]
		if !ok {
			fmt.Printf("%s is not a valid command\n", args[0])
			continue
		}
		err := cmd.callback(&config, args[1:])
		if err != nil {
			fmt.Println(err)
		}
	}
}
