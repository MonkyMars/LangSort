package scan

import (
	"filesorting/parse"
	"filesorting/structs"
	"fmt"
	"os"
	"path/filepath"
)

func ScanForFileSortFiles(rootDir string) ([]structs.FileSortConfig, error) {
	var configs []structs.FileSortConfig

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if this is a .filesort file
		if info.Name() == ".filesort" {
			config, err := parse.ParseFileSortFile(path)
			// Include the dir of the .filesort file in the config
			config.Dir = filepath.Dir(path)
			if config.Type == "" {
				fmt.Printf("Warning: .filesort file %s has no type defined, skipping.\n", path)
				return nil // Skip this file if no type is defined
			}
			if err != nil {
				fmt.Printf("Warning: Failed to parse %s: %v\n", path, err)
				return nil // Continue walking even if one file fails
			}
			configs = append(configs, config)
		}

		return nil
	})

	return configs, err
}

func ScanDirectSubdirectories(rootDir string) ([]structs.FileSortConfig, error) {
	var configs []structs.FileSortConfig

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if this directory has a .filesort file
		filesortPath := filepath.Join(rootDir, entry.Name(), ".filesort")
		if _, err := os.Stat(filesortPath); err == nil {
			config, err := parse.ParseFileSortFile(filesortPath)
			if err != nil {
				fmt.Printf("Warning: Failed to parse %s: %v\n", filesortPath, err)
				continue
			}
			// Include the dir of the .filesort file in the config
			config.Dir = filepath.Join(rootDir, entry.Name())
			if config.Type == "" {
				fmt.Printf("Warning: .filesort file %s has no type defined, skipping.\n", filesortPath)
				continue
			}
			configs = append(configs, config)
		}
	}

	return configs, nil
}
