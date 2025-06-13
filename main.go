package main

import (
	"filesorting/move"
	"filesorting/sanitize"
	"filesorting/scan"
	"fmt"
	"os"
	"path/filepath"
)

var RootDir = "/home/levinoppers/Coding"

func FolderExists(lang string) bool {
	// Expand the root directory path
	rootDir, err := filepath.Abs(RootDir)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return false
	}

	// Construct the full path to the language folder
	folderPath := filepath.Join(rootDir, lang)

	// Check if the folder exists
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func main() {
	fmt.Println("Scanning for .filesort files...")

	// Option 1: Recursive scan (finds .filesort files in all subdirectories)
	// configs, err := scan.ScanForFileSortFiles(RootDir)

	// Option 2: Only scan direct subdirectories (1 level deep)
	configs, err := scan.ScanDirectSubdirectories(RootDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	fmt.Printf("Found %d .filesort files:\n\n", len(configs))

	for _, config := range configs {
		fmt.Printf("  Type: %s\n", config.Type)
		fmt.Printf("  Directory: %s\n", config.Dir)
		fmt.Printf("\n")

		// Get the project name from the source directory
		projectName := filepath.Base(config.Dir)
		
		// Construct destination: /Coding/Language/ProjectName
		destDir := filepath.Join(RootDir, sanitize.Sanitize(config.Type), projectName)
		
		err := move.MoveDir(config.Dir, destDir, config.Type)
		if err != nil {
			fmt.Printf("Error moving %s: %v\n", config.Dir, err)
		}
	}
}
