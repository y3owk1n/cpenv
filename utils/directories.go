package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Directory struct {
	Name  string
	Value string
}

func IsFsDirectory(sourcePath string) (bool, error) {
	logrus.Debugf("Checking if path is a directory: %s", sourcePath)
	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("Path does not exist: %s", sourcePath)
			return false, fmt.Errorf("path does not exist: %s", sourcePath)
		}
		logrus.Errorf("Failed to stat path %s: %v", sourcePath, err)
		return false, fmt.Errorf("failed to stat path %s: %w", sourcePath, err)
	}
	logrus.Debugf("Path %s is a directory: %t", sourcePath, info.IsDir())
	return info.IsDir(), nil
}

func GetDirectories(directory string) ([]Directory, error) {
	logrus.Debugf("Getting directories in: %s", directory)
	entries, err := os.ReadDir(directory)
	if err != nil {
		logrus.Errorf("Failed to read directory %s: %v", directory, err)
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	directories := []Directory{}
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, Directory{
				Name:  entry.Name(),
				Value: entry.Name(),
			})
			logrus.Debugf("Found directory: %s", entry.Name())
		}
	}
	logrus.Debugf("Total directories found in %s: %d", directory, len(directories))
	return directories, nil
}

func Mkdir(path string) error {
	logrus.Debugf("Creating directory: %s", path)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		logrus.Errorf("Failed to create directory %s: %v", path, err)
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	logrus.Debugf("Directory created: %s", path)
	return nil
}

// It defaults to os.Getwd but can be overridden in tests.
var GetWdFunc = os.Getwd

func GetCurrentWorkingDirectory() string {
	dir, err := GetWdFunc()
	if err != nil {
		logrus.Errorf("Failed to get current working directory: %v", err)
		return ""
	}
	logrus.Debugf("Current working directory: %s", dir)
	return dir
}

func ReadDirRecursive(dirPath string) ([]string, error) {
	logrus.Debugf("Reading directory recursively: %s", dirPath)
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Errorf("Error accessing path %s: %v", path, err)
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("Failed to read directory recursively: %v", err)
		return nil, fmt.Errorf("failed to read directory recursively: %w", err)
	}
	logrus.Debugf("Total files found recursively in %s: %d", dirPath, len(files))
	return files, nil
}

// Id defaults to ReadDirRecursive but can be overridden in tests.
var ReadDirRecursiveFunc = ReadDirRecursive

func OpenInFinder(dirPath string) error {
	logrus.Debugf("Attempting to open directory in Finder: %s", dirPath)
	if runtime.GOOS != "darwin" {
		logrus.Errorf("OpenInFinder is only supported on macOS. Current OS: %s", runtime.GOOS)
		return fmt.Errorf("this function is only supported on macOS")
	}

	cmd := exec.Command("open", dirPath)
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("Failed to open directory %s in Finder: %v", dirPath, err)
		return fmt.Errorf("failed to open directory %s in Finder: %w", dirPath, err)
	}
	logrus.Debugf("Directory opened in Finder successfully: %s", dirPath)
	return nil
}
