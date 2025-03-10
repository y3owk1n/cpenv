package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createFakeEditor creates a fake editor executable that exits with the provided exitCode.
// On Unix-like systems, it creates a shell script, and on Windows, a batch file.
func createFakeEditor(t *testing.T, exitCode int) string {
	tempDir := t.TempDir()
	var scriptContent, fileName string
	if runtime.GOOS == "windows" {
		fileName = "fake_editor.bat"
		scriptContent = fmt.Sprintf("@echo off\r\nexit /b %d\r\n", exitCode)
	} else {
		fileName = "fake_editor"
		scriptContent = fmt.Sprintf("#!/bin/sh\nexit %d\n", exitCode)
	}
	fakeEditorPath := filepath.Join(tempDir, fileName)
	err := os.WriteFile(fakeEditorPath, []byte(scriptContent), 0755)
	assert.NoError(t, err)
	return fakeEditorPath
}

// TestOpenInEditorEmptyPath verifies that providing an empty file path returns an error.
func TestOpenInEditorEmptyPath(t *testing.T) {
	err := OpenInEditor("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file path cannot be empty")
}

// TestOpenInEditorFileNotExist verifies that calling OpenInEditor on a non-existent file returns an error.
func TestOpenInEditorFileNotExist(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.txt")
	err := OpenInEditor(nonExistentPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}

// TestOpenInEditorDefaultEditor verifies the branch where the EDITOR environment variable is not set.
// To avoid hanging (since it defaults to "vim"), we override PATH to an empty string so that "vim" cannot be found.
func TestOpenInEditorDefaultEditor(t *testing.T) {
	oldEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", oldEditor)
	os.Unsetenv("EDITOR")

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", "")

	tempFile, err := os.CreateTemp(t.TempDir(), "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = OpenInEditor(tempFile.Name())
	assert.Error(t, err)
	// The error from exec.Command when "vim" is not found usually mentions "executable file not found"
	assert.Contains(t, err.Error(), "executable file not found")
}

// TestOpenInEditorFakeEditorSuccess verifies that OpenInEditor succeeds when the fake editor exits with code 0.
func TestOpenInEditorFakeEditorSuccess(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "testfile")
	assert.NoError(t, err)
	// Close the file so Windows doesn't lock it
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	oldEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", oldEditor)
	fakeEditor := createFakeEditor(t, 0)
	os.Setenv("EDITOR", fakeEditor)

	err = OpenInEditor(tempFile.Name())
	assert.NoError(t, err)
}

// TestOpenInEditorFakeEditorFailure verifies that OpenInEditor returns an error when the fake editor exits with a non-zero code.
func TestOpenInEditorFakeEditorFailure(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	oldEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", oldEditor)
	fakeEditor := createFakeEditor(t, 1)
	os.Setenv("EDITOR", fakeEditor)

	err = OpenInEditor(tempFile.Name())
	assert.Error(t, err)
	// The fake editor exits with code 1, so the error message should contain "exit status 1"
	assert.Contains(t, err.Error(), "exit status 1")
}

// TestOpenInEditorStatError simulates a scenario where os.Stat returns an error that is not os.IsNotExist,
// such as a permission error. This test is skipped on Windows due to differences in permission handling.
func TestOpenInEditorStatError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping stat error test on Windows due to permission handling differences")
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	assert.NoError(t, err)

	// Remove permissions from the directory to force a permission error when calling os.Stat.
	origInfo, err := os.Stat(tempDir)
	assert.NoError(t, err)
	err = os.Chmod(tempDir, 0000)
	assert.NoError(t, err)
	// Restore permissions after the test.
	defer os.Chmod(tempDir, origInfo.Mode())

	err = OpenInEditor(filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking file")
}
