// Package utils implements the basic helper/utility functions to be used for
// github.com/maniac-en/pokefetch/cmd/pokefetch
package utils

import "strings"

func CleanInput(text string) []string {
	words := strings.Fields(text)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}
