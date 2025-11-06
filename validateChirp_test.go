package main

import (
	"testing"
)

func TestProfaneCheck(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello there",
			expected: "hello there",
		},
		{
			input:    " got this",
			expected: " got this",
		},
		{
			input:    "kerfuffle you",
			expected: "**** you",
		},
		{
			input:    "FORNAX me!",
			expected: "**** me!",
		},
		{
			input:    "sharbert!",
			expected: "sharbert!",
		},
	}

	for _, c := range cases {
		actual := profaneCheck(c.input)
		if c.expected != actual {
			t.Errorf("actual does not match expected, expected: %v, got: %v", c.expected, actual)
			t.Fail()
		}
	}
}
