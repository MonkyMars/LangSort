package config

import (
	"encoding/json"
	"filesorting/structs"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var Config = structs.Config{}

// expandPath expands ~ to the user's home directory and resolves relative paths
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %w", err)
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return filepath.Abs(path)
}

// findConfigFile looks for config.json in current directory, then in executable directory
func findConfigFile() (string, error) {
	configName := "config.json"

	// First, try current working directory
	if _, err := os.Stat(configName); err == nil {
		return configName, nil
	}

	// Then try the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, configName)

	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	return "", fmt.Errorf("config.json not found in current directory or executable directory")
}

func LoadConfig() error {
	configPath, err := findConfigFile()
	if err != nil {
		return err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close config file: %v\n", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	// Expand the directory path
	Config.Dir, err = expandPath(Config.Dir)
	if err != nil {
		return fmt.Errorf("failed to expand directory path: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

// validateConfig ensures the configuration is valid and the directory is accessible
func validateConfig() error {
	if Config.Dir == "" {
		return fmt.Errorf("sortDir cannot be empty")
	}

	if len(Config.AcceptedLanguages) == 0 {
		return fmt.Errorf("acceptedLanguages cannot be empty")
	}

	// Check if directory exists, create if it doesn't
	info, err := os.Stat(Config.Dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(Config.Dir, 0755); err != nil {
			return fmt.Errorf("failed to create sort directory %s: %w", Config.Dir, err)
		}
		fmt.Printf("Created sort directory: %s\n", Config.Dir)
	} else if err != nil {
		return fmt.Errorf("failed to check sort directory %s: %w", Config.Dir, err)
	} else if !info.IsDir() {
		return fmt.Errorf("sort path %s exists but is not a directory", Config.Dir)
	}

	// Test if directory is writable by creating a temporary file
	testFile := filepath.Join(Config.Dir, ".filesorting_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("sort directory %s is not writable: %w", Config.Dir, err)
	}
	file.Close()
	os.Remove(testFile) // Clean up test file

	return nil
}
