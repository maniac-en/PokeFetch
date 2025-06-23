package main

import (
	"errors"
	"fmt"
	"os"

	rand "math/rand/v2"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *string) error
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
		"explore": {
			name:        "explore",
			description: "Explore a map, like \"explore <map-area-name>\"",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pookemon",
			callback:    commandCatch,
		},
	}
}

func commandExit(cfg *config, _ *string) error {
	fmt.Println("Closing the PokeFetch... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, _ *string) error {
	fmt.Println("\nWelcome to the PokeFetch!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapf(cfg *config, _ *string) error {
	pokeMapAreas, err := cfg.client.GetMapAreas(cfg.nextMapAreaURL)
	if err != nil {
		return err
	}

	cfg.nextMapAreaURL = pokeMapAreas.Next
	cfg.prevMapAreaURL = pokeMapAreas.Previous

	for _, mapArea := range pokeMapAreas.Results {
		fmt.Println(mapArea.Name)
	}
	return nil
}

func commandMapb(cfg *config, _ *string) error {
	if cfg.prevMapAreaURL == nil {
		return errors.New("you're on the first page")
	}

	pokeMapAreas, err := cfg.client.GetMapAreas(cfg.prevMapAreaURL)
	if err != nil {
		return err
	}

	cfg.nextMapAreaURL = pokeMapAreas.Next
	cfg.prevMapAreaURL = pokeMapAreas.Previous

	for _, mapArea := range pokeMapAreas.Results {
		fmt.Println(mapArea.Name)
	}
	return nil
}

func commandExplore(cfg *config, param *string) error {
	if param == nil {
		return fmt.Errorf("can't explore empty map name, please provide a valid map name")
	}
	pokeMapArea, err := cfg.client.GetMapArea(param)
	if err != nil {
		return err
	}

	fmt.Println("\nFound Pokemon:")
	for _, pokemonEncounter := range pokeMapArea.PokemonEncounters {
		fmt.Println("-", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, param *string) error {
	if param == nil {
		return fmt.Errorf("can't catch a pokemon with no name, please provide one")
	}
	if pokemon, ok := cfg.pokedex[*param]; ok {
		fmt.Println("You already caught", pokemon.Name)
		return nil
	}
	pokemon, err := cfg.client.GetPokemon(param)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	chance := float64(rand.IntN(pokemon.BaseExperience))
	// fail if chance less than 40%
	if (chance / float64(pokemon.BaseExperience)) < 0.4 {
		fmt.Println(pokemon.Name, "escaped!")
	} else {
		fmt.Println(pokemon.Name, "was caught!")
		cfg.pokedex[pokemon.Name] = pokemon
	}
	return nil
}
