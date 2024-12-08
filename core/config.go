package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/utils"
)

func InitViper() error {
	viper.SetConfigName("cpenv")
	viper.SetConfigType("yaml")

	viper.SetDefault("vault_dir", ".env-files")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "cpenv")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}

	return nil
}

func GetFullVaultDir(vaultDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	vaultDirFull := filepath.Join(homeDir, vaultDir)

	utils.Logger.Debug("Vault path",
		"home_dir", homeDir,
		"vault_dir", vaultDir,
		"full_path", vaultDirFull,
	)

	return vaultDirFull, nil
}

func CreateVaultIfNotFound(vaultDir string) (string, error) {
	fullVaultDir, err := GetFullVaultDir(vaultDir)
	if err != nil {
		return "", fmt.Errorf("failed to get vault: %w", err)
	}

	_, err = os.Stat(fullVaultDir)

	if err == nil {
		utils.Logger.Debug("Vault already exists",
			"path", fullVaultDir,
		)
		fmt.Println(utils.WarningMessage.Render("Skipping creating vault, exists at", fullVaultDir))
		return fullVaultDir, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check vault: %w", err)
	}

	utils.Logger.Debug("Creating vault",
		"path", fullVaultDir,
	)

	if err := os.MkdirAll(fullVaultDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create vault: %w", err)
	}

	fmt.Println(utils.SuccessMessage.Render("Created vault directory at", fullVaultDir))

	return fullVaultDir, nil
}
