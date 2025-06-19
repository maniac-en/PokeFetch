package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	PROMPT string = "PokeFetch > "
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			break
		}

		inputLine := scanner.Text()
		cleanedInputLine := cleanInput(inputLine)
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

func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}
