package dbm

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a snake_case string to camelCase.
func ToCamelCase(s string, isUcFirst ...bool) string {
	var ucFirst bool
	if len(isUcFirst) > 0 {
		ucFirst = isUcFirst[0]
	}
	var camelCase string

	// Split the string into words separated by underscores.
	words := strings.Split(s, "_")

	// Capitalize the first letter of each word except the first one.
	for i, word := range words {
		if ucFirst || i > 0 {
			camelCase += strings.Title(word)
		} else {
			camelCase += word
		}
	}

	return camelCase
}

// ToSnakeCase converts a camelCase string to snake_case.
func ToSnakeCase(s string) string {
	var snakeCase string

	// Iterate over each character in the string.
	for i, r := range s {
		// If the character is uppercase and not the first character, add an underscore.
		if unicode.IsUpper(r) && i > 0 {
			snakeCase += "_"
		}
		// Convert the character to lowercase and append it to the result.
		snakeCase += strings.ToLower(string(r))
	}

	return snakeCase
}
