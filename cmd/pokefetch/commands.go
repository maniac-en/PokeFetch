package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"mapf": {
			name:        "mapf",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the PokeFetch... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("\nWelcome to the PokeFetch!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapf(cfg *config) error {
	pokeMaps, err := cfg.client.GetMapAreas(cfg.nextMapAreaURL)
	if err != nil {
		return err
	}

	cfg.nextMapAreaURL = pokeMaps.Next
	cfg.prevMapAreaURL = pokeMaps.Previous

	for _, mapArea := range pokeMaps.Results {
		fmt.Println(mapArea.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevMapAreaURL == nil {
		return errors.New("you're on the first page")
	}

	pokeMaps, err := cfg.client.GetMapAreas(cfg.prevMapAreaURL)
	if err != nil {
		return err
	}

	cfg.nextMapAreaURL = pokeMaps.Next
	cfg.prevMapAreaURL = pokeMaps.Previous

	for _, mapArea := range pokeMaps.Results {
		fmt.Println(mapArea.Name)
	}
	return nil
}
