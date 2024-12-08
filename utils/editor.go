package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenInEditor(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("error checking file %s: %w", filePath, err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
