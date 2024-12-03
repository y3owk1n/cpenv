package utils

import (
	"os"
	"os/exec"
)

func OpenInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	// Construct the command to open the file in the editor
	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and return any error
	return cmd.Run()
}
