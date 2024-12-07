package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Directory struct {
	Name  string
	Value string
}

// IsFsDirectory checks if the given path is a directory.
func IsFsDirectory(sourcePath string) (bool, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return false, fmt.Errorf("failed to stat path %s: %w", sourcePath, err)
	}
	return info.IsDir(), nil
}

// GetDirectories returns a list of directories within the specified path.
func GetDirectories(directory string) ([]Directory, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	directories := []Directory{}
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, Directory{
				Name:  entry.Name(),
				Value: entry.Name(),
			})
		}
	}
	return directories, nil
}

// Mkdir creates a directory at the specified path.
func Mkdir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func GetCurrentWorkingDirectory() string {
	// You can use os.Getwd() to get the current working directory
	dir, _ := os.Getwd()
	return dir
}

func ReadDirRecursive(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Only append files, not directories
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func OpenInFinder(dirPath string) error {
	// Check if the current OS is macOS
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this function is only supported on macOS")
	}

	// Use the 'open' command to open the directory in Finder
	cmd := exec.Command("open", dirPath)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open directory in Finder: %v", err)
	}

	return nil
}
