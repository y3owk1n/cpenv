package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type configInitCommand struct{}

type configEditCommand struct{}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Config management for cpenv",
	Aliases: []string{"c", "config"},
}

func newConfigInitCommand() *cobra.Command {
	cic := &configInitCommand{}

	return &cobra.Command{
		Use:     "init",
		Short:   "Initialize a config for cpenv to work",
		Aliases: []string{"i", "init"},
		Run:     cic.run,
	}
}

func newConfigEditCommand() *cobra.Command {
	cec := &configEditCommand{}

	return &cobra.Command{
		Use:              "edit",
		Short:            "Edit the cpenv config with $EDITOR",
		Aliases:          []string{"e", "edit"},
		PersistentPreRun: cec.preRun,
		Run:              cec.run,
	}
}

func (cic *configInitCommand) run(cmd *cobra.Command, args []string) {
	logrus.Debug("Starting config init command")

	if viper.ConfigFileUsed() != "" {
		fmt.Printf("%s %s\n", utils.ErrorIcon(), utils.WhiteText("Configuration exists! Use `cpenv config edit` to edit it"))
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Errorf("Failed to get user home directory: %v", err)
		return
	}
	logrus.Debugf("User home directory: %s", home)

	configPath := filepath.Join(home, ".config", "cpenv", "cpenv.yaml")
	logrus.Debugf("Config file path: %s", configPath)

	defaultVaultDir := ".env-files"
	logrus.Debugf("Default vault directory: %s", defaultVaultDir)

	viper.Set("vault_dir", defaultVaultDir)

	if err := viper.WriteConfigAs(configPath); err != nil {
		logrus.Errorf("Failed to save config: %v", err)
		return
	}
	logrus.Debug("Configuration written successfully")

	fmt.Printf("%s %s\n", utils.SuccessIcon(), utils.WhiteText("Configuration initialized successfully!"))

	if _, err := core.CreateVaultIfNotFound(defaultVaultDir); err != nil {
		logrus.Errorf("Failed to create vault directory: %v", err)
		return
	}
	logrus.Debug("Vault directory created successfully (if it was not already present)")
}

func (cec *configEditCommand) preRun(cmd *cobra.Command, args []string) {
	logrus.WithField("args", args).Debug("Starting config edit preRun command")

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Printf("%s %s\n", utils.ErrorIcon(), utils.WhiteText("Please run `cpenv config init` first"))
		os.Exit(0)
	}
	logrus.Debugf("Using config file: %s", configPath)

	vaultDir := viper.GetString("vault_dir")
	logrus.Debugf("Vault directory from config: %s", vaultDir)

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get env file directory: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Resolved full vault directory: %s", vaultDirFull)

	ctx := cmd.Context()
	ctx = context.WithValue(ctx, ConfigKey, configPath)
	ctx = context.WithValue(ctx, VaultKey, vaultDirFull)
	cmd.SetContext(ctx)
	logrus.Debugf("Context set with ConfigKey=%s and VaultKey=%s", configPath, vaultDirFull)
}

func (cec *configEditCommand) run(cmd *cobra.Command, args []string) {
	logrus.WithField("args", args).Debug("Starting config edit run command")

	configPath, ok := cmd.Context().Value(ConfigKey).(string)
	if !ok {
		logrus.Error("Config not found in context")
		fmt.Printf("%s %s\n", utils.ErrorIcon(), utils.WhiteText("vault config not found in context"))
	}
	logrus.Debugf("Retrieved config file from context: %s", configPath)

	if err := utils.OpenInEditor(configPath); err != nil {
		logrus.Errorf("Failed to open configuration file in editor: %v", err)
		os.Exit(1)
	}
	logrus.Debug("Configuration file opened in editor successfully")

	fmt.Printf("%s %s\n", utils.SuccessIcon(), utils.WhiteText("Successfully opened the configuration file in editor."))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to reload config: %v", err)
		os.Exit(1)
	}
	logrus.Debug("Configuration reloaded successfully")

	vaultDir := viper.GetString("vault_dir")
	logrus.Debugf("Vault directory after reload: %s", vaultDir)

	if _, err := core.CreateVaultIfNotFound(vaultDir); err != nil {
		logrus.Errorf("Failed to create vault directory: %v", err)
		os.Exit(1)
	}
	logrus.Debug("Vault directory ensured to exist after config reload")
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(newConfigInitCommand())
	configCmd.AddCommand(newConfigEditCommand())
}
