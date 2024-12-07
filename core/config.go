package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/y3owk1n/cpenv/utils"
)

type Config struct {
	VaultDir string `json:"vaultDir"`
}

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".env-files.json")

func LoadConfig() (*Config, error) {
	configType := reflect.TypeOf(Config{})
	configTypeName := fmt.Sprintf("%T", Config{})

	utils.Logger.Debug("Config type information",
		"reflect_type", configType,
		"type_name", configTypeName,
		"num_fields", configType.NumField(),
	)

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found at %s: %w", ConfigPath, err)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid config format: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", ConfigPath, err)
	}

	utils.Logger.Info("Configuration saved successfully",
		"path", ConfigPath,
		"vault_dir", config.VaultDir,
	)

	return nil
}

func validateConfig(config *Config) error {
	if config.VaultDir == "" {
		return fmt.Errorf("vault directory cannot be empty")
	}
	return nil
}

func GetEnvFilesDirectory(vaultDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	envFilesDirectory := filepath.Join(homeDir, vaultDir)

	utils.Logger.Debug("Env files directory path",
		"home_dir", homeDir,
		"vault_dir", vaultDir,
		"full_path", envFilesDirectory,
	)

	return envFilesDirectory, nil
}

func CreateEnvFilesDirectoryIfNotFound(vaultDir string) (string, error) {
	envFilesDirectory, err := GetEnvFilesDirectory(vaultDir)
	if err != nil {
		return "", fmt.Errorf("failed to get env files directory: %w", err)
	}

	_, err = os.Stat(envFilesDirectory)

	if err == nil {
		utils.Logger.Debug("Env files directory already exists",
			"path", envFilesDirectory,
		)
		return envFilesDirectory, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check env files directory: %w", err)
	}

	utils.Logger.Info("Creating env files directory",
		"path", envFilesDirectory,
	)

	if err := os.MkdirAll(envFilesDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create env files directory: %w", err)
	}

	return envFilesDirectory, nil
}

func TypeInformation() map[string]interface{} {
	configType := reflect.TypeOf(Config{})

	typeInfo := map[string]interface{}{
		"name":       configType.Name(),
		"kind":       configType.Kind(),
		"num_fields": configType.NumField(),
		"fields":     make([]map[string]string, configType.NumField()),
	}

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		typeInfo["fields"].([]map[string]string)[i] = map[string]string{
			"name":     field.Name,
			"type":     field.Type.String(),
			"json_tag": field.Tag.Get("json"),
		}
	}

	return typeInfo
}
