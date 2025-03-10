package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func OpenInEditor(filePath string) error {
	logrus.Debugf("Attempting to open file in editor: %s", filePath)
	if filePath == "" {
		logrus.Error("File path is empty")
		return fmt.Errorf("file path cannot be empty")
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			logrus.Errorf("File does not exist: %s", filePath)
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		logrus.Errorf("Error checking file %s: %v", filePath, err)
		return fmt.Errorf("error checking file %s: %w", filePath, err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		logrus.Debug("EDITOR environment variable not set; defaulting to vim")
		editor = "vim"
	} else {
		logrus.Debugf("Using editor from EDITOR environment variable: %s", editor)
	}

	logrus.Debugf("Executing command: %s %s", editor, filePath)
	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logrus.Errorf("Failed to open file %s in editor: %v", filePath, err)
	} else {
		logrus.Debugf("Successfully opened file %s in editor", filePath)
	}
	return err
}
