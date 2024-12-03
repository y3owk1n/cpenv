package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

func CheckFileExists(dir string, fileName string) (bool, error) {
	_, err := os.Stat(filepath.Join(dir, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CopyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func GetBackupTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}
