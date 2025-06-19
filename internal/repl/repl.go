// Package repl implements a basic Read-Eval-Print Loop for interacting with the PokeAPI.
// It handles user input and provides interactive commands in the terminal.
package repl

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/maniac-en/pokefetch/internal/client"
	"github.com/maniac-en/pokefetch/internal/utils"
)

const (
	PROMPT string = "PokeFetch > "
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
	config      struct {
		next     string
		previous string
	}
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config: struct {
				next     string
				previous string
			}{},
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config: struct {
				next     string
				previous string
			}{},
		},
		"map": {
			name:        "map",
			description: "Fetch next PokeMaps, 20 at a time",
			callback:    commandMap,
			config: struct {
				next     string
				previous string
			}{},
		},
		"mapb": {
			name:        "mapb",
			description: "Fetch previous PokeMaps, 20 at a time",
			callback:    commandMapb,
			config: struct {
				next     string
				previous string
			}{},
		},
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

var pokeMaps = &client.PokeMaps{}

func commandMap() error {
	if reflect.DeepEqual(pokeMaps, client.PokeMaps{}) {
		if err := client.GetMaps("", pokeMaps); err != nil {
			return err
		}
	} else {
		if err := client.GetMaps(pokeMaps.Next, pokeMaps); err != nil {
			return err
		}
	}
	list, err := printableMapNames(*pokeMaps)
	if err != nil {
		return err
	}
	fmt.Print(list)
	return nil
}

func commandMapb() error {
	if reflect.DeepEqual(pokeMaps, client.PokeMaps{}) {
		if err := client.GetMaps("", pokeMaps); err != nil {
			return err
		}
	} else {
		if err := client.GetMaps(*pokeMaps.Previous, pokeMaps); err != nil {
			return err
		}
	}
	list, err := printableMapNames(*pokeMaps)
	if err != nil {
		return err
	}
	fmt.Print(list)
	return nil
}

func printableMapNames(pokeMaps client.PokeMaps) (string, error) {
	if len(pokeMaps.Results) == 0 {
		return "", fmt.Errorf("got empty results from the API")
	}
	var list string
	for _, result := range pokeMaps.Results {
		list += fmt.Sprintf("%s\n", result.Name)
	}
	return list, nil
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			break
		}

		inputLine := scanner.Text()
		cleanedInputLine := utils.CleanInput(inputLine)
		if len(cleanedInputLine) == 0 {
			continue
		}
		inputCmd := cleanedInputLine[0]
		if handler, ok := commands[inputCmd]; !ok {
			fmt.Println("Unknown command:", inputCmd)
		} else {
			err := handler.callback()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: ", err)
	}
}
