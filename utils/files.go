package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CheckFileExists(dir string, fileName string) (bool, error) {
	if dir == "" || fileName == "" {
		return false, fmt.Errorf("directory and file name must not be empty")
	}

	filePath := filepath.Join(dir, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking file %s: %w", filePath, err)
	}
	return true, nil
}

func CopyFile(source, destination string) error {
	if source == "" || destination == "" {
		return fmt.Errorf("source and destination paths must not be empty")
	}

	srcFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", source, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destination, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file %s to %s: %w", destFile.Name(), srcFile.Name(), err)
	}

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy contents from %s to %s: %w", source, destination, err)
	}

	return nil
}

func GetBackupTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}
