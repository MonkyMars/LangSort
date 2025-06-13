package sanitize

import (
	"unicode"
)

// SanitizeString removes leading and trailing whitespace, replaces multiple spaces with a single space, and converts to lowercase.
func Sanitize(input string) string {
	var result []rune
	inSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !inSpace {
				result = append(result, ' ')
				inSpace = true
			}
		} else {
			result = append(result, unicode.ToLower(r))
			inSpace = false
		}
	}

	// Trim leading and trailing spaces
	if len(result) > 0 && result[0] == ' ' {
		result = result[1:]
	}
	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}

	return string(result)
}
