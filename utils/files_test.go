package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Tests for CheckFileExists ---

func TestCheckFileExistsEmptyInputs(t *testing.T) {
	// Both directory and fileName empty.
	exists, err := CheckFileExists("", "")
	assert.False(t, exists)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory and file name must not be empty")

	// Only directory empty.
	exists, err = CheckFileExists("", "file.txt")
	assert.False(t, exists)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory and file name must not be empty")

	// Only fileName empty.
	exists, err = CheckFileExists("someDir", "")
	assert.False(t, exists)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory and file name must not be empty")
}

func TestCheckFileExistsNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	exists, err := CheckFileExists(tempDir, "nonexistent.txt")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCheckFileExistsExists(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("hello"), 0644)
	assert.NoError(t, err)

	exists, err := CheckFileExists(tempDir, "test.txt")
	assert.NoError(t, err)
	assert.True(t, exists)
}

// TestCheckFileExistsStatError simulates a scenario where os.Stat returns an error
// that is not os.IsNotExist (e.g. permission error). This test is skipped on Windows.
func TestCheckFileExistsStatError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping stat error test on Windows due to permission handling differences")
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	assert.NoError(t, err)

	// Remove read permission from the directory to force a permission error.
	origInfo, err := os.Stat(tempDir)
	assert.NoError(t, err)
	err = os.Chmod(tempDir, 0000)
	assert.NoError(t, err)
	// Restore permissions after the test.
	defer os.Chmod(tempDir, origInfo.Mode())

	exists, err := CheckFileExists(tempDir, "test.txt")
	assert.False(t, exists)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking file")
}

// --- Tests for CopyFile ---

func TestCopyFileEmptyPaths(t *testing.T) {
	err := CopyFile("", "dest.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source and destination paths must not be empty")

	err = CopyFile("source.txt", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source and destination paths must not be empty")
}

func TestCopyFileNonExistentSource(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "nonexistent.txt")
	dest := filepath.Join(tempDir, "dest.txt")

	err := CopyFile(source, dest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open source file")
}

func TestCopyFileDestCreationError(t *testing.T) {
	tempDir := t.TempDir()
	// Create a valid source file.
	source := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(source, []byte("content"), 0644)
	assert.NoError(t, err)

	// Create a destination directory so that os.Create will fail.
	destDir := filepath.Join(tempDir, "destDir")
	err = os.Mkdir(destDir, 0755)
	assert.NoError(t, err)

	// Use the directory path as the destination file.
	err = CopyFile(source, destDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create destination file")
}

func TestCopyFileSuccess(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source.txt")
	originalContent := "Hello, world!"
	err := os.WriteFile(source, []byte(originalContent), 0644)
	assert.NoError(t, err)

	dest := filepath.Join(tempDir, "dest.txt")
	err = CopyFile(source, dest)
	assert.NoError(t, err)

	// Verify that the destination file was created and contains the same content.
	destContent, err := os.ReadFile(dest)
	assert.NoError(t, err)
	// Even though CopyFile calls io.Copy twice, the second copy reads zero bytes.
	assert.Equal(t, originalContent, string(destContent))
}

// --- Test for GetBackupTimestamp ---

func TestGetBackupTimestamp(t *testing.T) {
	timestamp := GetBackupTimestamp()
	// Parse the timestamp using the same format, ensuring local timezone.
	parsed, err := time.ParseInLocation("2006-01-02_15-04-05", timestamp, time.Local)
	assert.NoError(t, err)

	// Ensure the parsed time is reasonably close to now.
	now := time.Now()
	diff := now.Sub(parsed)
	if diff < 0 {
		diff = -diff
	}
	// Allow a margin of 2 seconds.
	assert.LessOrEqual(t, diff.Seconds(), 2.0, fmt.Sprintf("Timestamp difference too high: %v", diff))
}
