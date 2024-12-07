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

func IsFsDirectory(sourcePath string) (bool, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("path does not exist: %s", sourcePath)
		}
		return false, fmt.Errorf("failed to stat path %s: %w", sourcePath, err)
	}
	return info.IsDir(), nil
}

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

func Mkdir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func GetCurrentWorkingDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

func ReadDirRecursive(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read directory recursively: %w", err)
	}
	return files, nil
}

func OpenInFinder(dirPath string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this function is only supported on macOS")
	}

	cmd := exec.Command("open", dirPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open directory %s in Finder: %w", dirPath, err)
	}

	return nil
}
