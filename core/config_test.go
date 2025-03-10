package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestInitViperConfigNotFound tests the branch where no config file exists.
func TestInitViperConfigNotFound(t *testing.T) {
	// Save and restore HOME and reset viper.
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	viper.Reset()

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	// Ensure no config file exists.
	configDir := filepath.Join(tempHome, ".config", "cpenv")
	os.RemoveAll(configDir)

	err := InitViper()
	assert.NoError(t, err)

	// Verify that the config directory was created.
	_, err = os.Stat(configDir)
	assert.NoError(t, err)
}

// TestInitViperConfigInvalid tests the branch where a config file exists but is invalid.
func TestInitViperConfigInvalid(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	viper.Reset()

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	configDir := filepath.Join(tempHome, ".config", "cpenv")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Write an invalid YAML file.
	configFile := filepath.Join(configDir, "cpenv.yaml")
	err = os.WriteFile(configFile, []byte("invalid: [unbalanced brackets"), 0644)
	assert.NoError(t, err)

	err = InitViper()
	assert.Error(t, err)
}

// TestInitViperUserHomeError simulates an error from UserHomeDirFunc.
func TestInitViperUserHomeError(t *testing.T) {
	original := UserHomeDirFunc
	defer func() { UserHomeDirFunc = original }()

	UserHomeDirFunc = func() (string, error) {
		return "", fmt.Errorf("mocked error")
	}

	err := InitViper()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mocked error")
}

// TestGetFullVaultDir verifies that GetFullVaultDir returns the correct full path.
func TestGetFullVaultDir(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	vaultDir := "myVault"
	expected := filepath.Join(tempHome, vaultDir)

	result, err := GetFullVaultDir(vaultDir)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestGetFullVaultDirUserHomeError simulates an error from UserHomeDirFunc in GetFullVaultDir.
func TestGetFullVaultDirUserHomeError(t *testing.T) {
	original := UserHomeDirFunc
	defer func() { UserHomeDirFunc = original }()

	UserHomeDirFunc = func() (string, error) {
		return "", fmt.Errorf("mocked error")
	}

	result, err := GetFullVaultDir("vault")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get home directory")
	assert.Empty(t, result)
}

// TestCreateVaultIfNotFound_AlreadyExists verifies that if the vault directory already exists,
// CreateVaultIfNotFound returns its path without re-creating it.
func TestCreateVaultIfNotFound_AlreadyExists(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	vaultDir := ".env-files"
	fullVault := filepath.Join(tempHome, vaultDir)
	// Pre-create the vault directory.
	err := os.MkdirAll(fullVault, 0755)
	assert.NoError(t, err)

	result, err := CreateVaultIfNotFound(vaultDir)
	assert.NoError(t, err)
	assert.Equal(t, fullVault, result)
}

// TestCreateVaultIfNotFound_NotExist verifies that if the vault directory does not exist,
// CreateVaultIfNotFound creates it.
func TestCreateVaultIfNotFound_NotExist(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	vaultDir := "myVault"
	fullVault := filepath.Join(tempHome, vaultDir)
	os.RemoveAll(fullVault)

	result, err := CreateVaultIfNotFound(vaultDir)
	assert.NoError(t, err)
	assert.Equal(t, fullVault, result)

	// Check that the vault directory now exists.
	_, err = os.Stat(fullVault)
	assert.NoError(t, err)
}

// TestCreateVaultIfNotFound_MkdirAllError simulates a failure in creating the vault directory.
func TestCreateVaultIfNotFound_MkdirAllError(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	// Create a file where a directory is expected.
	parentPath := filepath.Join(tempHome, "parent")
	err := os.WriteFile(parentPath, []byte("not a directory"), 0644)
	assert.NoError(t, err)

	// Use a vaultDir that forces creation inside the "parent" file.
	vaultDir := filepath.Join("parent", "child")
	result, err := CreateVaultIfNotFound(vaultDir)
	assert.Error(t, err)
	assert.Empty(t, result)
}

// TestCreateVaultIfNotFoundUserHomeError simulates an error from UserHomeDirFunc in CreateVaultIfNotFound.
func TestCreateVaultIfNotFoundUserHomeError(t *testing.T) {
	original := UserHomeDirFunc
	defer func() { UserHomeDirFunc = original }()

	UserHomeDirFunc = func() (string, error) {
		return "", fmt.Errorf("mocked error")
	}

	result, err := CreateVaultIfNotFound("vault")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get vault")
	assert.Empty(t, result)
}
