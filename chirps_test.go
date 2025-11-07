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
			input:    " got  this-kerfuffle",
			expected: " got  this-kerfuffle",
		},
		{
			input:    "kerfuffle you",
			expected: "**** you",
		},
		{
			input:    "FORNAX me! sharbert",
			expected: "**** me! ****",
		},
		{
			input:    "sHArbert!",
			expected: "sHArbert!",
		},
		{
			input:    "kerfuffle fornax sharbert",
			expected: "**** **** ****",
		},
		{
			input:    "(sharbert)",
			expected: "(sharbert)",
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
