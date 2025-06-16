package main

import (
	"filesorting/config"
	"filesorting/move"
	"filesorting/sanitize"
	"filesorting/scan"
	"fmt"
	"os"
	"path/filepath"
)

func FolderExists(subPath string) bool {
	var targetPath string

	if subPath == "" {
		// Check if the base sort directory exists
		targetPath = config.Config.Dir
	} else {
		// Check if a subdirectory within the sort directory exists
		targetPath = filepath.Join(config.Config.Dir, subPath)
	}

	// Get absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		fmt.Printf("Error getting absolute path for %s: %v\n", targetPath, err)
		return false
	}

	// Check if the folder exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func main() {
	fmt.Println("Loading config...")
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Validate that the sort directory exists
	if !FolderExists("") { // Check if base sort directory exists
		fmt.Printf("Error: Sort directory %s does not exist\n", config.Config.Dir)
		return
	}

	fmt.Println("Scanning for .filesort files...")

	// Option 1: Recursive scan (finds .filesort files in all subdirectories)
	// configs, err := scan.ScanForFileSortFiles(RootDir)

	// Option 2: Only scan direct subdirectories (1 level deep)
	filesortConfigs, err := scan.ScanDirectSubdirectories(config.Config.Dir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	fmt.Printf("Found %d .filesort files:\n\n", len(filesortConfigs))

	for _, filesortConfig := range filesortConfigs {
		fmt.Printf("  Type: %s\n", filesortConfig.Type)
		fmt.Printf("  Directory: %s\n", filesortConfig.Dir)
		fmt.Printf("\n")

		// Get the project name from the source directory
		projectName := filepath.Base(filesortConfig.Dir)

		// Construct destination: /Coding/Language/ProjectName e.g.,
		destDir := filepath.Join(config.Config.Dir, sanitize.Sanitize(filesortConfig.Type), projectName)

		err := move.MoveDir(filesortConfig.Dir, destDir, filesortConfig.Type)
		if err != nil {
			fmt.Printf("Error moving %s: %v\n", filesortConfig.Dir, err)
		}
	}
}
