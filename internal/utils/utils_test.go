package utils

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		desc     string
		input    string
		expected []string
	}{
		{
			desc:     "basic single word",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			desc:     "multiple words with single spaces",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			desc:     "multiple words with leading/trailing spaces and multiple internal spaces",
			input:    "  hello   world  ",
			expected: []string{"hello", "world"},
		},
		{
			desc:     "single word with varied casing (should be lowercased)",
			input:    " hEllOwoRld ",
			expected: []string{"helloworld"},
		},
		{
			desc:     "mixed casing words (all should be lowercased)",
			input:    " HeLlO wOrLd ",
			expected: []string{"hello", "world"},
		},
		{
			desc:     "input only with spaces",
			input:    " ",
			expected: []string{},
		},
		{
			desc:     "input only with tabs",
			input:    "\t\t",
			expected: []string{},
		},
		{
			desc:     "input with mixed whitespace (spaces, tabs, newlines)",
			input:    "  word1\tword2\nword3  ", // strings.Fields splits by any Unicode space
			expected: []string{"word1", "word2", "word3"},
		},
		{
			desc:     "empty input string",
			input:    "",
			expected: []string{},
		},
		{
			desc:     "input with special characters (should be kept as is after lowercasing)",
			input:    "PokE-fETch!",
			expected: []string{"poke-fetch!"},
		},
		{
			desc:     "input with numbers",
			input:    "c0mmand 123",
			expected: []string{"c0mmand", "123"},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			actual := CleanInput(c.input)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("For input '%s': Expected %v, got %v", c.input, c.expected, actual)
			}
		})
	}
}
