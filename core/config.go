package core

import (
	"cpenv/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	VaultDir string `json:"vaultDir"`
}

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".env-files.json")

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found: %w", err)
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

func CreateEnvFilesDirectoryIfNotFound(vaultDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	envFilesDirectory := filepath.Join(homeDir, vaultDir)

	// Check if the directory exists
	_, err = os.Stat(envFilesDirectory)
	if err == nil {
		utils.Logger.Debugf("Env files directory exists at %s\n", envFilesDirectory)
		return envFilesDirectory, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check env files directory: %w", err)
	}

	utils.Logger.Debugf("Env files directory not found. Creating a new one at %s\n", envFilesDirectory)
	err = os.MkdirAll(envFilesDirectory, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create env files directory: %w", err)
	}

	utils.Logger.Debugf("Env files directory created at %s\n", envFilesDirectory)
	return envFilesDirectory, nil
}
