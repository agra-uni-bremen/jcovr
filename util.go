package main

import (
	"strings"
)

func increment(n int) int {
	return n + 1
}

func getLines(input string) []string {
	if len(input) == 0 {
		return []string{""} // empty file
	}

	// Remove terminating newline (if any)
	if input[len(input)-1] == '\n' {
		input = input[0 : len(input)-1]
	}

	return strings.Split(input, "\n")
}
