package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		test_desc string
		input     string
		expected  []string
	}{
		{
			test_desc: "basic single word",
			input:     "hello",
			expected:  []string{"hello"},
		},
		{
			test_desc: "multiple words with spaces around",
			input:     " hello  world ",
			expected:  []string{"hello", "world"},
		},
		{
			test_desc: "single word with spaces and varied casing",
			input:     " hEllOwoRld ",
			expected:  []string{"helloworld"},
		},
		{
			test_desc: "input with only spaces",
			input:     " ",
			expected:  []string{},
		},
		{
			test_desc: "empty input",
			input:     "",
			expected:  []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for i := range c.expected {
			if len(actual) != len(c.expected) {
				t.Log(c.test_desc)
				t.Errorf("Expected: %s, got %s", c.expected, actual)
				continue
			}
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Log(c.test_desc)
				t.Errorf("Expected: %s, got %s", expectedWord, word)
			}
		}
	}
}
