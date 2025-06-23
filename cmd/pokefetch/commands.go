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
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List out all the caught pokemons",
			callback:    commandPokedex,
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
		fmt.Println("You may now inspect it with the inspect command")
		cfg.pokedex[pokemon.Name] = pokemon
	}
	return nil
}

func commandInspect(cfg *config, param *string) error {
	if param == nil {
		return fmt.Errorf("can't inspect a pokemon with no name, please provide one")
	}
	if pokemon, ok := cfg.pokedex[*param]; ok {
		fmt.Println("Name:", pokemon.Name)
		fmt.Println("Height:", pokemon.Height)
		fmt.Println("Weight:", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
		return nil
	} else {
		return fmt.Errorf("you have not caught that pokemon")
	}
}

func commandPokedex(cfg *config, _ *string) error {
	if len(cfg.pokedex) == 0 {
		return fmt.Errorf("your pokedex is empty, go catch some pokemons with catch command")
	}
	fmt.Println("Your Pokedex:")
	for k := range cfg.pokedex {
		fmt.Println("  -", k)
	}
	return nil
}
