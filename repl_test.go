package main

import (
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
			desc:     "multiple words with spaces around",
			input:    " hello  world ",
			expected: []string{"hello", "world"},
		},
		{
			desc:     "single word with spaces and varied casing",
			input:    " hEllOwoRld ",
			expected: []string{"helloworld"},
		},
		{
			desc:     "input only with spaces",
			input:    " ",
			expected: []string{},
		},
		{
			desc:     "input only with tabs",
			input:    "		",
			expected: []string{},
		},
		{
			desc:     "empty input",
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for i := range c.expected {
			if len(actual) != len(c.expected) {
				t.Log(c.desc)
				t.Errorf("Expected: %s, got %s", c.expected, actual)
				continue
			}
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Log(c.desc)
				t.Errorf("Expected: %s, got %s", expectedWord, word)
			}
		}
	}
}
