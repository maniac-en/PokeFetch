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

func main() {
	for {
		fmt.Print(PROMPT)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputLine := scanner.Text()
			cleanedInputLine := cleanInput(inputLine)
			if len(cleanedInputLine) == 0 {
				fmt.Fprintln(os.Stderr, "Error: bad input")
				break
			}
			cmdToPrint := cleanedInputLine[0]
			fmt.Printf("Your command was: %s\n", cmdToPrint)
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input: ", err)
		}
	}
}

func cleanInput(text string) []string {
	output := []string{}
	splits := strings.SplitSeq(text, " ")
	for split := range splits {
		lowercaseSplit := strings.ToLower(split)
		trimmedSplit := strings.TrimSpace(lowercaseSplit)
		// ignore empty splits
		if trimmedSplit == "" {
			continue
		}
		output = append(output, trimmedSplit)
	}
	return output
}
