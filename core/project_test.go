package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"github.com/y3owk1n/cpenv/utils"
)

// ---------------------------
// Tests for GetProjectsList
// ---------------------------

func TestGetProjectsList_Error(t *testing.T) {
	// Use a non-existent vault directory so that utils.GetDirectories returns an error.
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
	projects, err := GetProjectsList(nonExistentDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error retrieving directories")
	assert.Nil(t, projects)
}

func TestGetProjectsList_NoProjects(t *testing.T) {
	// Create an empty temporary directory.
	tempDir := t.TempDir()
	projects, err := GetProjectsList(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no projects found in the vault")
	assert.Nil(t, projects)
}

func TestGetProjectsList_Success(t *testing.T) {
	// Create a temporary directory with two subdirectories.
	tempDir := t.TempDir()
	subDir1 := filepath.Join(tempDir, "proj1")
	subDir2 := filepath.Join(tempDir, "proj2")
	assert.NoError(t, os.Mkdir(subDir1, 0755))
	assert.NoError(t, os.Mkdir(subDir2, 0755))
	projects, err := GetProjectsList(tempDir)
	assert.NoError(t, err)
	// Verify that the project names match.
	var names []string
	for _, p := range projects {
		names = append(names, p.Name)
	}
	assert.Contains(t, names, "proj1")
	assert.Contains(t, names, "proj2")
}

// ---------------------------
// Tests for SelectProject
// ---------------------------

func TestSelectProject_EmptyProjects(t *testing.T) {
	selected, err := SelectProject([]utils.Directory{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no projects found in the vault")
	assert.Empty(t, selected)
}

func TestSelectProject_Success(t *testing.T) {
	// Override selectProjectRun to simulate a successful selection.
	origSelectRun := selectProjectRun
	defer func() { selectProjectRun = origSelectRun }()
	selectProjectRun = func(prompt promptui.Select) (int, string, error) {
		return 0, "project1", nil
	}
	projects := []utils.Directory{{Name: "project1"}, {Name: "project2"}}
	selected, err := SelectProject(projects)
	assert.NoError(t, err)
	assert.Equal(t, "project1", selected)
}

func TestSelectProject_NoSelection(t *testing.T) {
	// Simulate a prompt run that returns an empty string.
	origSelectRun := selectProjectRun
	defer func() { selectProjectRun = origSelectRun }()
	selectProjectRun = func(prompt promptui.Select) (int, string, error) {
		return 0, "", nil
	}
	projects := []utils.Directory{{Name: "project1"}}
	selected, err := SelectProject(projects)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no project selected")
	assert.Empty(t, selected)
}

func TestSelectProject_RunError(t *testing.T) {
	// Simulate a prompt run error.
	origSelectRun := selectProjectRun
	defer func() { selectProjectRun = origSelectRun }()
	testErr := fmt.Errorf("prompt error")
	selectProjectRun = func(prompt promptui.Select) (int, string, error) {
		return 0, "", testErr
	}
	projects := []utils.Directory{{Name: "project1"}}
	selected, err := SelectProject(projects)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error starting the selection form")
	assert.Empty(t, selected)
}

func TestSelectProject_PromptInterrupt(t *testing.T) {
	// Simulate a user interrupt.
	origSelectRun := selectProjectRun
	origExit := exitFunc
	defer func() {
		selectProjectRun = origSelectRun
		exitFunc = origExit
	}()
	selectProjectRun = func(prompt promptui.Select) (int, string, error) {
		return 0, "", promptui.ErrInterrupt
	}
	var exitCode int
	exitFunc = func(code int) {
		exitCode = code
		panic("exit")
	}
	projects := []utils.Directory{{Name: "project1"}}
	assert.Panics(t, func() { SelectProject(projects) })
	assert.Equal(t, 0, exitCode)
}

// ---------------------------
// Tests for generateProjectOptions
// ---------------------------

func TestGenerateProjectOptions(t *testing.T) {
	projects := []utils.Directory{
		{Name: "projA"},
		{Name: "projB"},
	}
	options := generateProjectOptions(projects)
	assert.Equal(t, []string{"projA", "projB"}, options)
}

// ---------------------------
// Tests for CopyEnvFilesToProject and related functions
// ---------------------------

func resolveTempDir(t *testing.T) string {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func TestPrettifiedPath(t *testing.T) {
	// Create a temporary base directory.
	cwdBaseDir := resolveTempDir(t)
	vaultBaseDir := resolveTempDir(t)

	// Change working directory to simulatedCwd so that os.Getwd() returns it.
	origWd, err := utils.GetWdFunc()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Fatal(err)
		}
	}()
	if err := os.Chdir(cwdBaseDir); err != nil {
		t.Fatal(err)
	}
	// Now os.Getwd() returns simulatedCwd.

	// Define test cases.
	// Note:
	// - When a file is under the simulated CWD, your function returns the unchanged absolute path.
	// - When under the vault directory, it returns a path with "{vault}/" as the prefix.
	// - When both CWD and vault are the same, your function currently returns "{vault}".
	// - Trailing slashes are not stripped in your function, so we expect them to be cleaned via filepath.Clean.
	tests := []struct {
		name     string
		path     string
		vaultDir string
		expected string
	}{
		// Updated macOS-compatible tests
		{
			name:     "CWD_Match",
			path:     filepath.Join(cwdBaseDir, "file.txt"),
			vaultDir: vaultBaseDir,
			expected: filepath.Join("{project}", "file.txt"),
		},
		{
			name:     "Vault_Match",
			path:     filepath.Join(vaultBaseDir, "data.db"),
			vaultDir: vaultBaseDir,
			expected: filepath.Join("{vault}", "data.db"),
		},
		{
			name:     "No_Match_Absolute",
			path:     "/etc/passwd",
			vaultDir: vaultBaseDir,
			expected: "/etc/passwd", // Works on Linux, macOS needs special handling
		},
		{
			name:     "No_Match_Relative_Outside",
			path:     filepath.Join("..", "outside.txt"),
			vaultDir: vaultBaseDir,
			expected: filepath.Join("..", "outside.txt"),
		},
	}

	// Run test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prettifiedPath(tt.path, tt.vaultDir)
			if got != tt.expected {
				t.Errorf("prettifiedPath(%q, %q) = %q; want %q", tt.path, tt.vaultDir, got, tt.expected)
			}
		})
	}
}

func TestCopyEnvFilesToProject_ReadDirError(t *testing.T) {
	// Provide a vaultDir that does not exist so that ReadDirRecursive returns an error.
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
	err := CopyEnvFilesToProject("dummyProject", "dummyCurrent", nonExistentDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading project path")
}

// normalizePath removes the "/private" prefix from paths on macOS.
func normalizePath(p string) string {
	if strings.HasPrefix(p, "/private/var/") {
		return "/var/" + strings.TrimPrefix(p, "/private/var/")
	}
	return p
}

func TestCopyEnvFilesToProject_Success(t *testing.T) {
	// Create a temporary directory structure for vaultDir.
	tempDir := t.TempDir()
	project := "proj"
	currentPath := "subdir"
	projectPath := filepath.Join(tempDir, project, currentPath)
	assert.NoError(t, os.MkdirAll(projectPath, 0755))

	// Create a dummy file inside the project.
	dummyFile := filepath.Join(projectPath, "file.env")
	content := []byte("test")
	assert.NoError(t, os.WriteFile(dummyFile, content, 0644))

	// Override copyFileWithSpinnerFunc to record its parameters.
	var recordedSource, recordedDest string
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		recordedSource = sourcePath
		recordedDest = destinationPath
		return nil
	}

	// Set the current working directory to a temporary directory.
	tempCwd := t.TempDir()
	origWd, err := utils.GetWdFunc()
	assert.NoError(t, err)
	err = os.Chdir(tempCwd)
	assert.NoError(t, err)
	defer func() {
		err = os.Chdir(origWd)
		assert.NoError(t, err)
	}()

	// Call the function under test.
	err = CopyEnvFilesToProject(project, currentPath, tempDir)
	assert.NoError(t, err)

	// Expected destination is based on the temporary working directory.
	expectedDest := filepath.Join(tempCwd, currentPath, "file.env")
	expectedNorm := normalizePath(expectedDest)
	recordedNorm := normalizePath(recordedDest)
	assert.Equal(t, expectedNorm, recordedNorm)
	assert.Equal(t, dummyFile, recordedSource)
}

func TestProcessCopyEnvFileToProject_FileNotExists(t *testing.T) {
	// Test branch where the destination file does not exist.
	tempProject := t.TempDir()
	currentPath := "current"
	dummyFile := filepath.Join(tempProject, "file.env")
	assert.NoError(t, os.WriteFile(dummyFile, []byte("data"), 0644))
	tempCwd := t.TempDir()
	origWd, _ := utils.GetWdFunc()
	os.Chdir(tempCwd)
	defer os.Chdir(origWd)
	destPath := filepath.Join(tempCwd, currentPath, "file.env")
	os.Remove(destPath)
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}
	err := processCopyEnvFileToProject(dummyFile, tempProject, currentPath, tempProject)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestProcessCopyEnvFileToProject_FileExists(t *testing.T) {
	// Create a temporary project directory.
	tempProject := t.TempDir()
	currentPath := "current"
	dummyFile := filepath.Join(tempProject, "file.env")
	assert.NoError(t, os.WriteFile(dummyFile, []byte("data"), 0644))
	// Create a temporary working directory to simulate destination.
	tempCwd := t.TempDir()
	destDir := filepath.Join(tempCwd, currentPath)
	assert.NoError(t, os.MkdirAll(destDir, 0755))
	destFile := filepath.Join(destDir, "file.env")
	assert.NoError(t, os.WriteFile(destFile, []byte("existing"), 0644))
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}
	origWd, err := utils.GetWdFunc()
	assert.NoError(t, err)
	err = os.Chdir(tempCwd)
	assert.NoError(t, err)
	defer func() {
		err = os.Chdir(origWd)
		assert.NoError(t, err)
	}()
	// Simulate user input "n\n" (do not overwrite) via os.Pipe.
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, err = w.WriteString("n\n")
	assert.NoError(t, err)
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()
	err = processCopyEnvFileToProject(dummyFile, tempProject, currentPath, tempProject)
	assert.NoError(t, err)
	// Since the destination file exists and user chose not to overwrite,
	// copyFileWithSpinnerFunc should not be called.
	assert.False(t, called, "Expected copy function not to be called when user declines overwrite")
}

func TestCopyFileWithSpinner(t *testing.T) {
	// Test copyFileWithSpinner with a temporary source file.
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	content := "hello"
	assert.NoError(t, os.WriteFile(sourceFile, []byte(content), 0644))
	destFile := filepath.Join(tempDir, "dest.txt")
	origCopyFile := utils.CopyFileFunc
	defer func() { utils.CopyFileFunc = origCopyFile }()
	utils.CopyFileFunc = func(source, destination string) error {
		data, err := os.ReadFile(source)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, data, 0644)
	}
	err := copyFileWithSpinner(sourceFile, destFile, tempDir)
	assert.NoError(t, err)
	data, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestHandleExistingFile_NotOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	// Simulate user input "n" so that the file is not overwritten.
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, err = w.WriteString("n\n")
	assert.NoError(t, err)
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}
	err = handleExistingFile("dummySource", "dummyDest", tempDir)
	assert.NoError(t, err)
	assert.False(t, called)
}

func TestHandleExistingFile_Overwrite(t *testing.T) {
	tempDir := t.TempDir()
	// Simulate user input "y" so that the file is overwritten.
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, err = w.WriteString("y\n")
	assert.NoError(t, err)
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}
	err = handleExistingFile("dummySource", "dummyDest", tempDir)
	assert.NoError(t, err)
	assert.True(t, called, "Expected copy function to be called when user confirms overwrite")
}

// ---------------------------
// Tests for ConfirmCwd
// ---------------------------

func TestConfirmCwd_Success(t *testing.T) {
	// Set working directory to a temporary directory.
	tempDir := t.TempDir()
	origWd, _ := utils.GetWdFunc()
	os.Chdir(tempDir)
	defer os.Chdir(origWd)
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, err = w.WriteString("y\n")
	assert.NoError(t, err)
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()
	err = ConfirmCwd()
	assert.NoError(t, err)
}

func TestConfirmCwd_Abort(t *testing.T) {
	origExit := exitFunc
	var exitCode int
	exitFunc = func(code int) {
		exitCode = code
		panic("exit")
	}
	defer func() { exitFunc = origExit }()
	tempDir := t.TempDir()
	origWd, _ := utils.GetWdFunc()
	os.Chdir(tempDir)
	defer os.Chdir(origWd)
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, err = w.WriteString("n\n")
	assert.NoError(t, err)
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()
	assert.Panics(t, func() { ConfirmCwd() })
	assert.Equal(t, 0, exitCode)
}

// ---------------------------
// Tests for CopyEnvFilesToVault and processCopyEnvFileToVault
// ---------------------------

func TestCopyEnvFilesToVault_Error(t *testing.T) {
	// Use a temporary vault directory to avoid creating a "vaultDir" in the real project.
	tempVault := t.TempDir()
	origReadDir := utils.ReadDirRecursiveFunc
	defer func() { utils.ReadDirRecursiveFunc = origReadDir }()
	utils.ReadDirRecursiveFunc = func(dirPath string) ([]string, error) {
		return nil, fmt.Errorf("failed to read directory recursively: simulated error")
	}
	err := CopyEnvFilesToVault(tempVault)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading project path")
}

func TestCopyEnvFilesToVault_WithTempVault(t *testing.T) {
	// Create a temporary directory to simulate the vault.
	tempVaultDir := t.TempDir()

	// Create a temporary directory to simulate the current project.
	tempProjectDir := t.TempDir()
	// Create a dummy file in the project directory.
	dummyFile := filepath.Join(tempProjectDir, "config.env")
	err := os.WriteFile(dummyFile, []byte("content"), 0644)
	assert.NoError(t, err)

	// Override utils.ReadDirRecursiveFunc to simulate reading the dummy file.
	origReadDir := utils.ReadDirRecursiveFunc
	defer func() { utils.ReadDirRecursiveFunc = origReadDir }()
	utils.ReadDirRecursiveFunc = func(dirPath string) ([]string, error) {
		// Assume that the current project directory is used.
		return []string{dummyFile}, nil
	}

	// Override copyFileWithSpinnerFunc to record that it was called.
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}

	// Simulate that the current working directory is the temporary project directory.
	origWd, err := utils.GetWdFunc()
	assert.NoError(t, err)
	err = os.Chdir(tempProjectDir)
	assert.NoError(t, err)
	defer func() {
		err = os.Chdir(origWd)
		assert.NoError(t, err)
	}()

	// Call CopyEnvFilesToVault with the temporary vault directory.
	err = CopyEnvFilesToVault(tempVaultDir)
	assert.NoError(t, err)
	assert.True(t, called, "Expected copyFileWithSpinnerFunc to be called")

	// Verify that a backup folder was created in the vault.
	projectBase := filepath.Base(tempProjectDir)
	entries, err := os.ReadDir(tempVaultDir)
	assert.NoError(t, err)
	found := false
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), projectBase+"-") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected a backup folder starting with the project name in the vault")
}

func TestProcessCopyEnvFileToVault_SkipCases(t *testing.T) {
	tempDir := t.TempDir()
	// Test branches that should skip copying.
	err := processCopyEnvFileToVault("node_modules/somefile.env", "dummyCwd", "dummyDest", tempDir)
	assert.NoError(t, err)
	err = processCopyEnvFileToVault("dummyCwd/config.template", "dummyCwd", "dummyDest", tempDir)
	assert.NoError(t, err)
	err = processCopyEnvFileToVault("dummyCwd/config.example", "dummyCwd", "dummyDest", tempDir)
	assert.NoError(t, err)
	err = processCopyEnvFileToVault("dummyCwd/readme.txt", "dummyCwd", "dummyDest", tempDir)
	assert.NoError(t, err)
}

func TestProcessCopyEnvFileToVault_Copy(t *testing.T) {
	tempDir := t.TempDir()
	// Test the branch that actually copies an .env file.
	tempCwd := t.TempDir()
	dummyFile := filepath.Join(tempCwd, "config.env")
	assert.NoError(t, os.WriteFile(dummyFile, []byte("content"), 0644))
	var called bool
	origCopySpinner := copyFileWithSpinnerFunc
	defer func() { copyFileWithSpinnerFunc = origCopySpinner }()
	copyFileWithSpinnerFunc = func(sourcePath, destinationPath string, vaultDir string) error {
		called = true
		return nil
	}
	err := processCopyEnvFileToVault(dummyFile, tempCwd, filepath.Join(tempCwd, "vault"), tempDir)
	assert.NoError(t, err)
	assert.True(t, called)
}
