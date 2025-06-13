package parse

import (
	"filesorting/structs"
	"fmt"
	"os"
	"strings"
)

func ParseFileSortFile(filepath string) (structs.FileSortConfig, error) {
	var config structs.FileSortConfig

	// Read the file content
	content, err := os.ReadFile(filepath)
	if err != nil {
		return config, fmt.Errorf("error reading file: %w", err)
	}

	// Split the content into lines
	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		// Split the line into key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return config, fmt.Errorf("invalid line format: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "type":
			config.Type = value
		default:
			return config, fmt.Errorf("unknown key: %s", key)
		}
	}
	return config, nil
}
