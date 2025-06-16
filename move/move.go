package move

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func folderExists(dir string) bool {
	// Clean the path for cross-platform compatibility
	dir = filepath.Clean(dir)

	// Get absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Printf("Error getting absolute path for %s: %v\n", dir, err)
		return false
	}

	// Check if the folder exists
	info, err := os.Stat(absDir)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Printf("Error checking directory %s: %v\n", absDir, err)
		return false
	}
	return info.IsDir()
}

func MoveDir(sourceDir, destDir, lang string) error {
	// Clean and normalize paths for cross-platform compatibility
	sourceDir = filepath.Clean(sourceDir)
	destDir = filepath.Clean(destDir)

	// Check if source directory exists
	if exists := folderExists(sourceDir); !exists {
		return fmt.Errorf("source directory %s does not exist", sourceDir)
	}

	// Create the parent directory of the destination if it doesn't exist
	parentDir := filepath.Dir(destDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("error creating parent directory %s: %w", parentDir, err)
	}

	// Check if destination already exists
	if folderExists(destDir) {
		return fmt.Errorf("destination directory %s already exists", destDir)
	}

	// Try to move the source directory to the destination folder using os.Rename
	// This should automatically remove the source directory
	if err := os.Rename(sourceDir, destDir); err != nil {
		// If os.Rename fails (e.g., across filesystems), fall back to copy + delete
		fmt.Printf("Direct rename failed, falling back to copy + delete: %v\n", err)

		if err := copyDir(sourceDir, destDir); err != nil {
			return fmt.Errorf("error copying directory from %s to %s: %w", sourceDir, destDir, err)
		}

		// On Windows, there might be a slight delay before the file system recognizes the copy is complete
		// Add a small retry mechanism for the removal
		if err := removeWithRetry(sourceDir, 3); err != nil {
			return fmt.Errorf("error removing source directory %s after copy: %w", sourceDir, err)
		}

		fmt.Printf("Successfully copied and removed %s to %s\n", sourceDir, destDir)
	} else {
		fmt.Printf("Successfully moved %s to %s\n", sourceDir, destDir)
	}

	// Double-check that source directory no longer exists
	if folderExists(sourceDir) {
		return fmt.Errorf("source directory %s still exists after move operation", sourceDir)
	}

	return nil
}

// removeWithRetry attempts to remove a directory with retries for Windows compatibility
func removeWithRetry(path string, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := os.RemoveAll(path); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				// Wait a bit before retrying (useful on Windows)
				time.Sleep(time.Millisecond * 100)
				continue
			}
		} else {
			return nil
		}
	}
	return lastErr
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from src
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Calculate the destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Printf("Warning: Failed to close source file: %v\n", err)
		}
	}()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			fmt.Printf("Warning: Failed to close destination file: %v\n", err)
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
