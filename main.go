package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	output := []string{}
	splits := strings.Split(text, " ")
	for _, split := range splits {
		// ignore empty splits
		if split == "" {
			continue
		}
		lowercase_split := strings.ToLower(split)
		trimmed_split := strings.Trim(lowercase_split, " ")
		output = append(output, trimmed_split)
	}
	return output
}
