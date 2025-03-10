package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsFsDirectory tests the IsFsDirectory function.
func TestIsFsDirectory(t *testing.T) {
	// Case 1: Path is a directory.
	dir := t.TempDir()
	isDir, err := IsFsDirectory(dir)
	assert.NoError(t, err)
	assert.True(t, isDir)

	// Case 2: Path is a file.
	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	isDir, err = IsFsDirectory(tmpFile.Name())
	assert.NoError(t, err)
	assert.False(t, isDir)

	// Case 3: Path does not exist.
	nonExist := filepath.Join(t.TempDir(), "nonexistent")
	_, err = IsFsDirectory(nonExist)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path does not exist")
}

// TestGetDirectories tests the GetDirectories function.
func TestGetDirectories(t *testing.T) {
	// Create a temporary base directory.
	baseDir := t.TempDir()
	// Create two subdirectories.
	subDir1 := filepath.Join(baseDir, "sub1")
	subDir2 := filepath.Join(baseDir, "sub2")
	err := os.MkdirAll(subDir1, 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(subDir2, 0755)
	assert.NoError(t, err)
	// Create a file inside the base directory.
	filePath := filepath.Join(baseDir, "file.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	assert.NoError(t, err)

	// Get directories and check.
	dirs, err := GetDirectories(baseDir)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dirs))
	expectedNames := map[string]bool{"sub1": true, "sub2": true}
	for _, d := range dirs {
		assert.True(t, expectedNames[d.Name], "unexpected directory: %s", d.Name)
	}

	// Error case: non-existent directory.
	_, err = GetDirectories(filepath.Join(baseDir, "nonexistent"))
	assert.Error(t, err)
}

// TestMkdir tests the Mkdir function.
func TestMkdir(t *testing.T) {
	baseDir := t.TempDir()
	newDir := filepath.Join(baseDir, "newDir")

	// Create a new directory.
	err := Mkdir(newDir)
	assert.NoError(t, err)
	info, err := os.Stat(newDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Creating the directory again should succeed.
	err = Mkdir(newDir)
	assert.NoError(t, err)

	// Error case: try to create a directory where a file exists.
	filePath := filepath.Join(baseDir, "aFile")
	err = os.WriteFile(filePath, []byte("content"), 0644)
	assert.NoError(t, err)
	err = Mkdir(filePath)
	assert.Error(t, err)
}

// TestGetCurrentWorkingDirectory tests GetCurrentWorkingDirectory.
func TestGetCurrentWorkingDirectory(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	result := GetCurrentWorkingDirectory()
	assert.Equal(t, cwd, result)
}

// TestReadDirRecursive tests the ReadDirRecursive function.
func TestReadDirRecursive(t *testing.T) {
	baseDir := t.TempDir()

	// Create a flat file.
	file1 := filepath.Join(baseDir, "file1.txt")
	err := os.WriteFile(file1, []byte("content1"), 0644)
	assert.NoError(t, err)

	// Create nested directories with files.
	subDir1 := filepath.Join(baseDir, "sub1")
	err = os.MkdirAll(subDir1, 0755)
	assert.NoError(t, err)
	file2 := filepath.Join(subDir1, "file2.txt")
	err = os.WriteFile(file2, []byte("content2"), 0644)
	assert.NoError(t, err)

	subSubDir := filepath.Join(subDir1, "subsub")
	err = os.MkdirAll(subSubDir, 0755)
	assert.NoError(t, err)
	file3 := filepath.Join(subSubDir, "file3.txt")
	err = os.WriteFile(file3, []byte("content3"), 0644)
	assert.NoError(t, err)

	// Create an empty directory.
	emptyDir := filepath.Join(baseDir, "empty")
	err = os.MkdirAll(emptyDir, 0755)
	assert.NoError(t, err)

	// Execute ReadDirRecursive.
	files, err := ReadDirRecursive(baseDir)
	assert.NoError(t, err)
	// Expect only files (not directories).
	expectedFiles := []string{file1, file2, file3}
	sort.Strings(files)
	sort.Strings(expectedFiles)
	assert.Equal(t, expectedFiles, files)

	// Error case: non-existent directory.
	_, err = ReadDirRecursive(filepath.Join(baseDir, "nonexistent"))
	assert.Error(t, err)
}

// TestOpenInFinder tests the OpenInFinder function.
func TestOpenInFinder(t *testing.T) {
	tempDir := t.TempDir()

	// For non-darwin systems, the function should immediately return an error.
	if runtime.GOOS != "darwin" {
		err := OpenInFinder(tempDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "this function is only supported on macOS")
		return
	}

	// On macOS, simulate the open command by overriding the PATH so that a fake "open" is used.
	// Save the original PATH and restore it at the end.
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Create a temporary directory to hold our fake "open" command.
	fakeDir := t.TempDir()

	// Helper to create a fake "open" command.
	createFakeOpen := func(exitCode int) error {
		fakeOpenPath := filepath.Join(fakeDir, "open")
		script := fmt.Sprintf("#!/bin/sh\nexit %d\n", exitCode)
		return os.WriteFile(fakeOpenPath, []byte(script), 0755)
	}

	// Prepend the fakeDir to PATH.
	newPath := fakeDir + string(os.PathListSeparator) + origPath
	assert.NoError(t, os.Setenv("PATH", newPath))

	// Test success: fake "open" returns 0.
	err := createFakeOpen(0)
	assert.NoError(t, err)
	err = OpenInFinder(tempDir)
	assert.NoError(t, err)

	// Test failure: fake "open" returns non-zero.
	err = createFakeOpen(1)
	assert.NoError(t, err)
	err = OpenInFinder(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open directory")
}
