package utils

import (
	"strings"
)

func CleanString(str string) string {
	// Replace all newline and tab characters with a space
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\t", " ")

	// Trim any leading and trailing white spaces
	return strings.TrimSpace(str)
}
