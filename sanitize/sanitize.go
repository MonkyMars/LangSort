package sanitize

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// invalidCharsRegex matches characters that are invalid in filenames on Windows
var invalidCharsRegex = regexp.MustCompile(`[<>:"/\\|?*]`)

// Sanitize removes leading and trailing whitespace, replaces multiple spaces with a single space,
// converts to lowercase, and removes invalid filename characters for cross-platform compatibility
func Sanitize(input string) string {
	if input == "" {
		return ""
	}

	var result []rune
	inSpace := false

	// First pass: normalize whitespace and convert to lowercase
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

	sanitized := string(result)

	// Remove invalid filename characters for cross-platform compatibility
	sanitized = invalidCharsRegex.ReplaceAllString(sanitized, "")

	// Replace any remaining problematic characters with underscores
	sanitized = strings.ReplaceAll(sanitized, " ", "_")

	// Ensure the result is safe for use as a directory name
	sanitized = filepath.Clean(sanitized)

	// Handle empty result or result that starts with a dot (hidden files)
	if sanitized == "" || sanitized == "." {
		sanitized = "unknown"
	}

	return sanitized
}
