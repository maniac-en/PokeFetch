package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/maniac-en/pokefetch/internal/client"
	"github.com/maniac-en/pokefetch/internal/utils"
)

type config struct {
	client         client.Client
	nextMapAreaURL *string
	prevMapAreaURL *string
	pokedex        map[string]client.Pokemon
}

const (
	PROMPT string = "PokeFetch > "
)

func ReplStart(cfg *config) {
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
		var params *string
		if len(cleanedInputLine) > 1 {
			restOfTheInput := strings.Join(cleanedInputLine[1:], " ")
			params = &restOfTheInput
		}
		if handler, ok := getCommands()[inputCmd]; !ok {
			fmt.Println("Unknown command:", inputCmd)
		} else {
			err := handler.callback(cfg, params)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: ", err)
	}
}
