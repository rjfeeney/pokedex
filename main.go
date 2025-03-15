package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, []string) error
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	cleaned := strings.Fields(lowered)
	return cleaned
}

func commandExit(_ *pokeapi.Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *pokeapi.Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for command, cmd := range commands {
		fmt.Printf("%s: %s\n", command, cmd.description)
	}
	return nil
}

var pokeClient *pokeapi.Client

func main() {
	pokeClient = pokeapi.NewClient()
	config := &pokeapi.Config{}
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback: func(config *pokeapi.Config, _ []string) error {
				return commandExit(config, nil)
			},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(config *pokeapi.Config, _ []string) error {
				return commandHelp(config, nil)
			},
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback: func(config *pokeapi.Config, _ []string) error {
				return pokeClient.CommandMap(config)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback: func(config *pokeapi.Config, _ []string) error {
				return pokeClient.CommandMapb(config)
			},
		},
		"explore": {
			name:        "explore",
			description: "Displays all pokemon found in the area",
			callback: func(config *pokeapi.Config, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("please enter a location area")
				}
				return pokeClient.CommandExplore(args[0])
			},
		},
		"catch": {
			name:        "catch",
			description: "Throw a Pokeball for a chance to catch the Pokemon",
			callback: func(config *pokeapi.Config, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("please enter a Pokemon")
				}
				return pokeClient.CommandCatch(args[0])
			},
		},
		"inspect": {
			name:        "inspect",
			description: "Displays the info for a Pokemon if it has been caught",
			callback: func(config *pokeapi.Config, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("please enter a Pokemon")
				}
				return pokeClient.CommandInspect(args[0])
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists Pokemon that have been caught",
			callback: func(config *pokeapi.Config, _ []string) error {
				return pokeClient.CommandPokedex()
			},
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			words := cleanInput(text)
			if len(words) > 0 {
				command := words[0]
				args := []string{}
				if len(words) > 1 {
					args = words[1:]
				}
				cmd, found := commands[command]
				if found {
					err := cmd.callback(config, args)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			} else {
				fmt.Println("No input provided.")
			}
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				break
			}
		}
	}
}
