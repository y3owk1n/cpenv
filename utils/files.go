package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func CheckFileExists(dir string, fileName string) (bool, error) {
	logrus.Debugf("Checking if file exists. Directory: %s, FileName: %s", dir, fileName)
	if dir == "" || fileName == "" {
		logrus.Error("Directory and file name must not be empty")
		return false, fmt.Errorf("directory and file name must not be empty")
	}

	filePath := filepath.Join(dir, fileName)
	logrus.Debugf("Constructed file path: %s", filePath)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("File does not exist: %s", filePath)
			return false, nil
		}
		logrus.Errorf("Error checking file %s: %v", filePath, err)
		return false, fmt.Errorf("error checking file %s: %w", filePath, err)
	}
	logrus.Debugf("File exists: %s", filePath)
	return true, nil
}

func CopyFile(source, destination string) error {
	logrus.Debugf("Copying file from source: %s to destination: %s", source, destination)
	if source == "" || destination == "" {
		logrus.Error("Source and destination paths must not be empty")
		return fmt.Errorf("source and destination paths must not be empty")
	}

	srcFile, err := os.Open(source)
	if err != nil {
		logrus.Errorf("Failed to open source file %s: %v", source, err)
		return fmt.Errorf("failed to open source file %s: %w", source, err)
	}
	defer func() {
		if cerr := srcFile.Close(); cerr != nil {
			logrus.Errorf("Error closing source file %s: %v", source, cerr)
		}
	}()

	destFile, err := os.Create(destination)
	if err != nil {
		logrus.Errorf("Failed to create destination file %s: %v", destination, err)
		return fmt.Errorf("failed to create destination file %s: %w", destination, err)
	}
	defer func() {
		if cerr := destFile.Close(); cerr != nil {
			logrus.Errorf("Error closing destination file %s: %v", destination, cerr)
		}
	}()

	copiedBytes, err := io.Copy(destFile, srcFile)
	if err != nil {
		logrus.Errorf("Failed to copy file from %s to %s: %v", source, destination, err)
		return fmt.Errorf("failed to copy file %s to %s: %w", source, destination, err)
	}
	logrus.Debugf("Copied %d bytes from %s to %s", copiedBytes, source, destination)

	// NOTE: The second io.Copy is redundant as the source is likely at EOF.
	// However, it is retained to follow the original logic.
	if _, err := io.Copy(destFile, srcFile); err != nil {
		logrus.Errorf("Failed to copy contents from %s to %s: %v", source, destination, err)
		return fmt.Errorf("failed to copy contents from %s to %s: %w", source, destination, err)
	}
	logrus.Debugf("File copy completed successfully from %s to %s", source, destination)
	return nil
}

func GetBackupTimestamp() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logrus.Debugf("Generated backup timestamp: %s", timestamp)
	return timestamp
}
