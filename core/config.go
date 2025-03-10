package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/utils"
)

func InitViper() error {
	logrus.Debug("Initializing Viper configuration")
	viper.SetConfigName("cpenv")
	viper.SetConfigType("yaml")

	viper.SetDefault("vault_dir", ".env-files")
	logrus.Debug("Set default vault_dir to .env-files")

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Errorf("Failed to get user home directory: %v", err)
		return err
	}
	logrus.Debugf("User home directory: %s", home)

	configPath := filepath.Join(home, ".config", "cpenv")
	logrus.Debugf("Config path: %s", configPath)

	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	if err := os.MkdirAll(configPath, 0755); err != nil {
		logrus.Errorf("Failed to create config directory: %v", err)
		return err
	}
	logrus.Debug("Config directory ensured to exist")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Debug("No config file found; continuing without it")
			return nil
		}
		logrus.Errorf("Failed to read config: %v", err)
		return err
	}

	logrus.Debug("Viper configuration loaded successfully")
	return nil
}

func GetFullVaultDir(vaultDir string) (string, error) {
	logrus.Debugf("Resolving full vault directory for vault_dir: %s", vaultDir)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Errorf("Failed to get home directory: %v", err)
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	logrus.Debugf("User home directory: %s", homeDir)

	vaultDirFull := filepath.Join(homeDir, vaultDir)
	logrus.Debugf("Full vault directory resolved: %s", vaultDirFull)

	return vaultDirFull, nil
}

func CreateVaultIfNotFound(vaultDir string) (string, error) {
	logrus.Debugf("Ensuring vault exists for vault_dir: %s", vaultDir)
	fullVaultDir, err := GetFullVaultDir(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get full vault directory: %v", err)
		return "", fmt.Errorf("failed to get vault: %w", err)
	}

	logrus.Debugf("Checking if vault directory exists: %s", fullVaultDir)
	_, err = os.Stat(fullVaultDir)
	if err == nil {
		logrus.Debugf("Vault already exists at: %s", fullVaultDir)
		return fullVaultDir, nil
	}

	if !os.IsNotExist(err) {
		logrus.Errorf("Error checking vault directory: %v", err)
		return "", fmt.Errorf("failed to check vault: %w", err)
	}

	logrus.Debugf("Vault does not exist. Creating vault directory: %s", fullVaultDir)
	if err := os.MkdirAll(fullVaultDir, 0755); err != nil {
		logrus.Errorf("Failed to create vault directory: %v", err)
		return "", fmt.Errorf("failed to create vault: %w", err)
	}

	fmt.Printf("%s %s %s\n", utils.SuccessIcon(), utils.WhiteText("Created vault directory at"), utils.CyanText(fullVaultDir))
	logrus.Debugf("Vault directory created at: %s", fullVaultDir)

	return fullVaultDir, nil
}
