package util

import (
	"fmt"
	"strings"
)

func Confirm(s string, defaultValue bool) bool {
	var input string

	fmt.Print(s)
	fmt.Print(" [Y/n]: ")
	fmt.Scanln(&input)

	if input == "" {
		return defaultValue
	}

	// Convert the input to lowercase for case-insensitive comparison
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	} else {
		return false
	}
}
